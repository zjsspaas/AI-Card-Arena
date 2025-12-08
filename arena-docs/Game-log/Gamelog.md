from dataclasses import dataclass, field, asdict
from typing import List, Any, Dict
import json


# ---------------------------------------------------------
#  日志数据结构（符合你的规范）
# ---------------------------------------------------------
@dataclass
class GameLog:
    game_id: str
    turn: int
    self_seat: int
    landlord: int

    current_hand: str
    others_hand: str

    num_cards_left: List[int]
    actions: List[str]

    played_cards: List[str] = field(default_factory=list)
    action_record: List[List[Any]] = field(default_factory=list)

    schema_version: str = "1.0.0"

    # ------------------ 日志操作函数 ------------------

    def add_action(self, player: int, action: str):
        """
        记录一次出牌动作，自动更新 turn
        """
        self.action_record.append([player, action])
        self.turn = len(self.action_record) + 1

        if action != "pass":
            self.played_cards.append(action)

    def to_json(self) -> str:
        """
        导出 JSON 字符串，复盘系统将加载这个
        """
        return json.dumps(asdict(self), ensure_ascii=False, indent=2)

    @staticmethod
    def from_json(json_str: str):
        """
        从 JSON 字符串加载日志（复盘）
        """
        data = json.loads(json_str)
        return GameLog(
            game_id=data["game_id"],
            turn=data["turn"],
            self_seat=data["self"],
            landlord=data["landlord"],
            current_hand=data["current_hand"],
            others_hand=data["others_hand"],
            num_cards_left=data["num_cards_left"],
            actions=data["actions"],
            played_cards=data.get("played_cards", []),
            action_record=data.get("action_record", []),
            schema_version=data.get("schema_version", "1.0.0")
        )

    # ---------------------------------------------------------
    # 复盘辅助：按动作编号回放状态
    # ---------------------------------------------------------
    def get_state_at(self, step: int) -> Dict[str, Any]:
        """
        获取 action_record 第 step 步之后的局面，
        用于前端复盘时展示某一帧的牌局信息。
        """
        step = max(0, min(step, len(self.action_record)))

        state = {
            "game_id": self.game_id,
            "turn": step,
            "action_record": self.action_record[:step],
            "played_cards": [],
        }

        # 取前 step 个非 pass 动作
        for p, act in self.action_record[:step]:
            if act != "pass":
                state["played_cards"].append(act)

        return state
/**************************************************/
Game-log
-这是一个 类（class），代表一局游戏（对局）的日志结构。
game_id
-游戏／对局的唯一标识符 (string)。用于区分不同对局，比如 "game136137"。
turn
-当前回合数 (int) —— 表示“当前是第几轮动作”，方便知道游戏进度。
self_seat
-“我” (agent) 的座位号 (int)，值为 0, 1, 或 2。对应你在输入合约里 self。
landlord
-地主玩家的座位号 (int) (0/1/2)。对应你输入合约里的 landlord。
current_hand
-我的当前手牌 (string)。按照你之前定义的编码规则编码。
others_hand
-其余所有未在我手上的牌 (string)。
num_cards_left
-每个玩家剩余手牌数量 (List[int])。
-格式是 [player0_count, player1_count, player2_count]。便于快速知道每人还有多少牌。
actions
-当前可行动作列表 (List[str])。
-表示当前回合，我（agent）可选择的合法动作 (比如 "pass", "77", "99" 等
played_cards
-到目前为止所有已打出的牌 (List[str])。
-每一项是一次出牌 (非 “pass”) 的字符串。用于复盘／展示已出的牌堆。
action_record
-完整动作记录 (List[List[Any]])。
-这是一个二维列表，每个子列表代表一条动作记录，格式为 [player_index, action_string]，按顺序记录了每一步谁出了什么 (或 pass)。
schema_version
-日志 schema 的版本号 (string)。默认为 "1.0.0"。方便未来如果你修改结构，可以用版本号来兼容 / 区分不同格式。
