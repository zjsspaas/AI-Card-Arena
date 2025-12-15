from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Dict, Any
import rlcard
import uuid
app = FastAPI(title="AI Card Arena - Python Engine")
# 初始化 FastAPI 应用实例，这是整个服务的入口

# 内存存储：保持所有活跃牌局的状态
game_envs: Dict[str, Any] = {}
# 【核心状态层】定义一个全局字典，用 game_id (str) 作为键，存储 RLCard 的游戏环境对象 (env)。
# 这是实现多局游戏并发的核心机制。

class StepAction(BaseModel):
# 定义 Pydantic 模型，用于校验 /internal/game/step 接口的请求体 JSON 格式
    game_id: str
    # 接收 Go 后端传来的游戏唯一标识符
    action: str # 对应动作：出牌或不出
    # 接收 Go 后端传来的 Agent 决策动作，如 ["K", "K"] 或 ["pass"]

def translate_to_contract(env, player_id: int):
    """
    [解释]：将算法库状态翻译为《Agent 接入合约》要求的格式
    """
    # 此函数是将 RLCard 的原始数据转化为 Agent 和 Go 后端能理解的标准化 JSON 格式的“翻译官”

    state = env.get_state(player_id)
    # 调用 RLCard 方法，获取当前玩家视角的完整游戏状态
    obs = state['raw_obs']
    # 提取 RLCard 状态中的原始观察数据 (raw_obs)，这里面通常是人类可读的字符串牌面

    return {
        "turn"=
        "landlord":env.game.round.landlord_id,
        # 判定当前玩家在牌局中的身份 (地主或农民)
        "hand_cards": obs['hand'], #
        # 获取当前玩家的手牌列表，例如 ["3", "4", "5", "RJ"]
        "public_info": {
            "last_move": obs['trace'][-1] if obs['trace'] else None,
            # 获取桌面上最后一手牌的信息 (牌型和出牌者 ID)
            "remaining_cards": {
                f"player_{i}": count for i, count in enumerate(obs['others_hand'])
            } #
            # 统计并返回其他玩家剩余的牌数，是 AI 决策的关键信息
        },
        "legal_actions": state['raw_legal_actions']
        # 获取当前玩家所有合法的出牌动作列表 (这是 Agent 决策的约束集合)
    }

@app.post("/internal/game/init")
# 定义 API 路由：这是 Go 后端用于启动新牌局的接口
async def init_game():
    """初始化接口 """
    game_id = str(uuid.uuid4())
    # 使用 uuid 库生成一个全局唯一的 ID，作为牌局的唯一标识
    env = rlcard.make('doudizhu') #
    # 初始化 RLCard 的斗地主游戏环境对象
    env.reset()
    # 执行牌局重置和初始化（即洗牌和发牌）

    game_envs[game_id] = env
    # 将创建好的 env 实例存储到全局内存字典中
    next_p = env.get_player_id()
    # 获取第一个行动的玩家 ID (通常是叫牌阶段的第一个人)

    return {
        "game_id": game_id,
        "next_player_id": next_p,
        "state": translate_to_contract(env, next_p)
        # 返回游戏 ID 和第一个玩家的初始状态
    }

@app.post("/internal/game/step")
# 定义 API 路由：这是 Go 后端用于执行动作和推进牌局的核心接口
async def process_step(req: StepAction):
    """推进牌局并执行动作"""
    # FastAPI 自动将请求体 JSON 验证并转换为 StepAction 对象
    if req.game_id not in game_envs:
        raise HTTPException(status_code=404, detail="Game not found")
    # 【鲁棒性检查】如果 Go 传来的 game_id 不存在，返回 404

    env = game_envs[req.game_id]
    # 根据 game_id 从内存中取出对应的游戏环境实例
    action_str = "".join(req.action) if req.action else "pass" #
    # 将 Agent 返回的动作列表（如 ["3", "3"]）拼接成 RLCard 要求的字符串格式（如 "33" 或 "pass"）

    try:
        _, next_p = env.step(action_str) # 推进规则逻辑
        # 【核心规则执行】将动作传入 RLCard，由其自动更新状态，并返回下一个玩家 ID
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Illegal Move: {str(e)}")
        # 【规则校验】如果动作非法（如出牌不符合牌型或比上家小），RLCard 抛出异常。
        # Python 服务捕获异常并返回 400，通知 Go 裁判该 Agent 违规。

    is_over = env.is_over()
    # 检查牌局是否结束
    response = {
        "game_id": req.game_id,
        "is_over": is_over,
        "next_player_id": next_p if not is_over else None
    }
    # 构造基本响应结构

    if is_over:
        response["reward"] = env.get_payoffs().tolist() #
        # 如果游戏结束，获取玩家的最终奖赏/惩罚 (用于结算和日志)
        del game_envs[req.game_id] # 释放内存
        # 【内存管理】游戏结束后，从内存中删除该 env 实例，防止内存泄漏
    else:
        response["state"] = translate_to_contract(env, next_p)
        # 如果游戏未结束，获取下一个玩家的最新状态，并准备返回给 Go

    return response
    # 返回最终结果 JSON 给 Go 后端

if __name__ == "__main__":
    import uvicorn
    # 仅在直接运行此文件时执行以下代码
    uvicorn.run(app, host="0.0.0.0", port=8000)
    # 启动 ASGI 服务器 (Uvicorn)，监听所有 IP 的 8000 端口，等待 Go 后端调用
