from dataclasses import dataclass
from enum import Enum
from typing import Dict, List, Optional, Any

import logging

from .engine import GameEngine

logger = logging.getLogger(__name__)

class RoomState(Enum):
    WAITING = "waiting"
    PLAYING = "playing"
    FINISHED = "finished"


@dataclass
class RoomConfig:
    seed: Optional[int] = None


class GameRoom:
    def __init__(self, room_id: str, config: RoomConfig):
        # 房间唯一 ID
        self.room_id = room_id
        self.config = config

        # 房间当前状态
        self.status = RoomState.WAITING

        # 一局对应一个引擎实例
        self.engine = GameEngine(config.seed)

        # 当前轮到的玩家
        self._current_player_id: Optional[int] = None

        # 地主是谁
        self._landlord_id: Optional[int] = None

        # 记录动作: [[player_id, action_str], ...]
        self.action_record: List[List[Any]] = []  # 例如 [[0, "55588"], [1, "JJQQQ"], [2, "pass"], ...]

        # 只记录非 pass 的出牌
        self.played_cards: List[str] = []

        # 是否已经初始化
        self._is_initialized = False

    def initialize(self) -> Dict[str, Any]:  # 返回首个玩家状态
        if self._is_initialized:
            raise ValueError("Room already initialized")

        # 让引擎重置一局游戏
        state, player_id = self.engine.reset()

        # 当前轮到谁
        self._current_player_id = player_id

        # 地主是谁
        self._landlord_id = self.engine.get_landlord_id()

        # 更新房间状态
        self.status = RoomState.PLAYING
        self._is_initialized = True

        logger.info(
            f"房间 {self.room_id} 初始化完成，地主是玩家 {self._landlord_id}，"
            f"当前轮到玩家 {self._current_player_id}"
        )
        return self.get_state_for_player(player_id)

    def get_state_for_player(self, player_id: int) -> Dict[str, Any]:
        if not self._is_initialized:
            raise ValueError("房间还没初始化，不能获取状态")

        # 从引擎获取 RLCard 状态
        rl_state = self.engine.get_state(player_id)
        raw_obs = rl_state.get('raw_obs', {})

        # 当前玩家手牌（简化：直接拼成字符串）
        hand = raw_obs.get('hand', [])
        current_hand = "".join(hand)

        # 合法动作列表（字符串）
        actions = self.engine.get_legal_actions(player_id)

        # 其他玩家剩余牌数
        others_left = raw_obs.get('others_hand', [])
        if len(others_left) != 2:
            others_left = [0, 0]

        # 三个玩家各自的剩余牌数
        num_cards_left = [0, 0, 0]
        num_cards_left[player_id] = len(hand)
        other_ids = [i for i in range(3) if i != player_id]
        num_cards_left[other_ids[0]] = others_left[0]
        num_cards_left[other_ids[1]] = others_left[1]

        # 拷贝一份出牌记录和动作记录，防止外部修改内部状态
        played_cards = self.played_cards[:]
        action_record = [step[:] for step in self.action_record]

        # 牌堆中剩下但不在自己手上的牌（这里先简化为空字符串）
        others_hand = ""

        # 当前是第几轮（简单用 action_record 长度 + 1）
        turn = len(self.action_record) + 1

        return {
            "game_id": self.room_id,
            "turn": turn,
            "self": player_id,
            "landlord": self._landlord_id,
            "current_hand": current_hand,
            "others_hand": others_hand,
            "num_cards_left": num_cards_left,
            "actions": actions,
            "played_cards": played_cards,
            "action_record": action_record,
        }

    def step(self, player_id: int, action: str) -> Dict[str, Any]:
        """
        执行一步出牌动作

        返回：
        - 如果游戏还没结束：
          {
            "status": "playing",
            "next_player": int,
            "state": {...下一玩家状态...}
          }

        - 如果游戏已经结束：
          {
            "status": "finished",
            "payoffs": [...],
            "winner": int,
            "action_record": [...]
          }
        """
        if not self._is_initialized:
            raise ValueError("房间还没初始化，不能出牌")

        if self.status != RoomState.PLAYING:
            raise ValueError(f"当前房间状态不是 PLAYING，而是 {self.status}")

        if player_id != self._current_player_id:
            raise ValueError(f"现在轮到玩家 {self._current_player_id} 出牌，不是玩家 {player_id}")

        # 检查动作是否合法
        legal_actions = self.engine.get_legal_actions(player_id)
        if action not in legal_actions:
            raise ValueError(f"非法动作: {action}，合法动作有: {legal_actions}")
        try:
            _, next_player_id = self.engine.step(action)
        except Exception as e:
            raise ValueError(f"执行动作失败: {e}")

        # 更新当前玩家
        self._current_player_id = next_player_id

        # 记录 [player_id, action]，方便日志与复盘
        self.action_record.append([player_id, action])

        # 非 pass 动作计入已出牌
        if action != "pass":
            self.played_cards.append(action)

        # 如果游戏已经结束
        if self.engine.is_over():
            self.status = RoomState.FINISHED
            payoffs = self.engine.get_payoffs()
            winner = self._get_winner(payoffs)
            return {
                "status": "finished",
                "payoffs": payoffs,
                "winner": winner,
                "action_record": [step[:] for step in self.action_record],
            }

        # 游戏未结束，返回下一玩家的状态
        next_state = self.get_state_for_player(next_player_id)
        return {
            "status": "playing",
            "next_player": next_player_id,
            "state": next_state,
        }

    def _get_winner(self, payoffs: List[float]) -> int:
        """根据 RLCard 的 payoffs 判断哪个玩家赢了"""
        landlord = self._landlord_id
        landlord_score = payoffs[landlord]

        if landlord_score > 0:
            return landlord

        # 农民胜利：从另外两个玩家中选得分更高的那位
        farmer_ids = [i for i in range(3) if i != landlord]
        farmer_scores = [(i, payoffs[i]) for i in farmer_ids]
        winner_id, _ = max(farmer_scores, key=lambda x: x[1])
        return winner_id

