from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Dict, Any, Optional
import uuid

from engine.game.room import GameRoom, RoomConfig

# 创建 FastAPI 应用
app = FastAPI(title="AI Card Arena - Python Engine")

# 用一个全局字典，在内存里保存所有正在进行的牌局
# key: game_id (string)，value: GameRoom 实例
rooms: Dict[str, GameRoom] = {}

class InitGameRequest(BaseModel):
    """初始化一局游戏时，Go 传过来的请求体"""
    seed: Optional[int] = None  # 随机种子，可以不传


class InitGameResponse(BaseModel):
    """初始化一局游戏时，Python 返回给 Go 的数据"""
    game_id: str
    state: Dict[str, Any]  # 直接把 GameRoom 返回的状态原样塞进去


class StepRequest(BaseModel):
    """推进一手牌的请求体"""
    game_id: str   # 要推进哪一局
    player_id: int # 哪个玩家在出牌（0/1/2）
    action: str    # 出的牌，比如 "pass" 或 "33TTKKKAAA"


class StepResponse(BaseModel):
    """推进一手牌后的响应体"""
    status: str                         # "playing" 或 "finished"
    next_player: Optional[int] = None   # 如果还在进行，告诉下一个玩家是谁
    state: Optional[Dict[str, Any]] = None  # 如果还在进行，下一个玩家的状态
    payoffs: Optional[list[float]] = None   # 如果结束了，返回最终得分
    winner: Optional[int] = None            # 如果结束了，返回赢家座位号

@app.post("/internal/game/init", response_model=InitGameResponse)
def init_game(req: InitGameRequest):
    """
    创建一局新的斗地主牌局，返回 game_id 和首个玩家的状态
    """
    try:
        # 1. 生成一个唯一的 game_id
        game_id = f"game_{uuid.uuid4().hex[:12]}"

        # 2. 创建 GameRoom 实例
        room = GameRoom(room_id=game_id, config=RoomConfig(seed=req.seed))

        # 3. 调用 room.initialize()，发牌并确定地主，得到第一个要出牌玩家的状态
        state = room.initialize()

        # 4. 把房间对象放到全局字典里，方便后续通过 game_id 找回
        rooms[game_id] = room

        # 5. 返回给 Go：game_id + 初始状态
        return InitGameResponse(game_id=game_id, state=state)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"初始化游戏失败: {str(e)}")


@app.post("/internal/game/step", response_model=StepResponse)
def step_game(req: StepRequest):
    """
    推进一手牌，返回下一步状态或游戏结束信息
    """
    try:
        # 1. 从全局字典中找到对应的房间
        room = rooms.get(req.game_id)
        if not room:
            raise HTTPException(status_code=404, detail=f"游戏 {req.game_id} 不存在")

        # 2. 执行出牌动作
        result = room.step(req.player_id, req.action)

        # 3. 返回结果
        if result["status"] == "finished":
            # 游戏结束，清理房间
            del rooms[req.game_id]
            return StepResponse(
                status="finished",
                payoffs=result["payoffs"],
                winner=result["winner"]
            )
        else:
            # 游戏继续
            return StepResponse(
                status="playing",
                next_player=result["next_player"],
                state=result["state"]
            )
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"执行动作失败: {str(e)}")


@app.get("/internal/game/{game_id}/state")
def get_game_state(game_id: str, player_id: int):
    """
    获取指定玩家的游戏状态（用于重连或查询）
    """
    room = rooms.get(game_id)
    if not room:
        raise HTTPException(status_code=404, detail=f"游戏 {game_id} 不存在")
    
    try:
        state = room.get_state_for_player(player_id)
        return {"game_id": game_id, "state": state}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.delete("/internal/game/{game_id}")
def delete_game(game_id: str):
    """
    删除游戏房间（用于清理或取消游戏）
    """
    if game_id in rooms:
        del rooms[game_id]
        return {"message": f"游戏 {game_id} 已删除"}
    raise HTTPException(status_code=404, detail=f"游戏 {game_id} 不存在")


@app.get("/health")
def health_check():
    """
    健康检查接口
    """
    return {
        "status": "healthy",
        "active_games": len(rooms),
        "service": "AI Card Arena - Python Engine"
    }