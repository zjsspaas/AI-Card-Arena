# agent须遵守的牌局信息输入合约
### 示例：
```json
{
  "game_id": "game136137",
  "turn": 12,
  "self": 2,
  "landlord": 0,
  "current_hand": "35677899TJK",
  "others_hand": "3778JQK",
  "num_cards_left": [6,1,11],
  "actions": ["pass", "77", "99"],
  "played_cards": ["45558899AAA22R","344666TTTJJQQQKB","34KA22"],
  "action_record": [
      [0, "55588"], [1, "JJQQQ"], [2, "pass"],
      [0, "AAA22"], [1, "pass"], [2, "pass"],
      [0, "9"], [1, "pass"], [2, "A"],
      [0, "pass"], [1, "B"], [2, "pass"],
      [0, "pass"], [1, "TT"], [2, "22"],
      [0, "pass"], [1, "pass"],[2, "4"],
      [0, "R"], [1, "pass"], [2, "pass"],
      [0, "4"], [1, "T"], [2, "K"],
      [0, "pass"], [1, "pass"],[2, "3"],
      [0, "9"],[1, "K"],[2, "pass"],
      [0, "pass"],[1, "3666"],[2, "pass"],
      [0, "pass"],[1, "44"]
    ]
}
```
* game_id:本场游戏唯一的id
* turn:当前是第几轮出牌
* self:你的座位号(0,1,2)
* landlord:地主座位号(0,1,2)
* current_hand：你的手牌（T为10，R为大王，B为小王）
* others_hand：剩余牌堆信息，表示牌堆中未被自己持有的牌，但不能直接确定对手的具体手牌。
* num_cards_left:所有玩家的剩余手牌数量[玩家0手牌数,玩家1手牌数,玩家2手牌数]
* actions：当前的合法动作： Pass（不出）、出对7、出对9。
* played_cards：所有玩家已打出的牌
* action_record：所有历史出牌动作