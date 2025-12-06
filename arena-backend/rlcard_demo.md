# 🃏 基于 RLCard 的斗地主随机智能体对战示例

本工程演示如何基于 **RLCard 开源框架** 构建一个可运行的斗地主环境，并让三个随机智能体在环境中自动对战。
代码完整复现了 RLCard 官方环境的核心流程，包括：

* 环境初始化
* 智能体注册
* 完整的对局循环
* 使用环境提供的 Perfect Information 打印初始手牌
* 动作 ID 解码到真实可读牌面
* 终局奖励显示

本项目可作为：
✔ RLCard 初学者示例
✔ 强化学习智能体项目的基础模板
✔ 探索斗地主动作空间与状态空间的参考代码

---

# 📂 项目结构

```
project/
│
├── rlcard_demo.py               # 主要示例程序（你的代码）
├── README.md             # 本说明文档
└── requirements.txt      # 项目依赖（RLCard 等）
```

---

# 🔧 环境准备

### 1. 安装依赖

```
pip install rlcard
pip install numpy
```
### 2. （可选）创建 requirements.txt
```
rlcard>=1.0.4
numpy
```

---

# 🚀 程序运行方式

```
python rlcard_demo.py
```

运行后你会看到：

* 地主/农民身份
* 初始发牌（带花色）
* 每步出牌动作的 Action-ID 与真实牌面字符串
* 游戏循环
* 最终胜负与 payoff

---

# 🧩 程序结构与 RLCard 对应源码解析
以下按顺序说明 main.py 中的核心部分，并标注 RLCard 源码关联位置，确保读者理解代码与 RLCard 框架的对应关系。

---

## 1. 环境初始化（与 RLCard Env 类对应）

```python
env = rlcard.make('doudizhu', config={'seed': 42})
set_seed(42)
```
### ✔ 对应 RLCard 源码位置
| 模块                        | 文件路径                    | 功能               |
| `rlcard/envs/doudizhu.py` | `DoudizhuEnv.__init__`     | 构建游戏实例、动作空间、玩家数量 |
| `rlcard/envs/env.py`      | `Env.reset()`、`Env.step()` | 通用环境逻辑           |
| `rlcard/games/doudizhu`   |                            | 斗地主游戏主逻辑         |
## 2. 构建随机智能体（与 RandomAgent 对应）
```python
for i in range(env.num_players):
    agent = RandomAgent(num_actions=env.num_actions)
```
斗地主动作空间约 300 个（包括 PASS、单牌、对子、炸弹、顺子等）。

---
## 3. 获取初始手牌（来自 Perfect Information）
```python
perfect_info = env.get_perfect_information()
hands = perfect_info['hand_cards_with_suit']
```
### ✔ RLCard 源码关联
| 位置                                      | 说明                |
| --------------------------------------- | ----------------- |
| `DoudizhuEnv.get_perfect_information()` | 返回完整牌局信息（用于评估/日志） |
| 字段：`hand_cards_with_suit`               | 所有玩家完整手牌（带花色）     |
---
## 4. 斗地主地主确定（Landlord ID）
```python
landlord_id = env.game.round.landlord_id
```
### ✔ 对应 RLCard 源码
| 模块                               | 路径                  | 说明     |
| -------------------------------- | ------------------- | ------ |
| `rlcard/games/doudizhu/round.py` | `Round.landlord_id` | 标记地主玩家 |
斗地主总是从地主开始出牌。
---
## 5. 游戏循环（标准 RLCard 机制）
```python
while not env.is_over():
    current_player_id = env.get_player_id()
    state = env.get_state(current_player_id)
    action = agents[current_player_id].step(state)
    next_state, next_player_id = env.step(action)
```
### ✔ 源码关联
| 位置                         | 功能          |
| -------------------------- | ----------- |
| `Env.get_state(player_id)` | 获取当前玩家的可见状态 |
| `RandomAgent.step(state)`  | 随机选 action  |
| `Env.step(action)`         | 推进状态、切换玩家   |
| `Game.is_over()`           | 判断是否结束      |
---
## 6. 动作 ID 解码（斗地主非常重要）
```python
action_card_str = env._ID_2_ACTION[action]
```
### ✔ RLCard 动作结构说明
在 `DoudizhuEnv.__init__` 中有字段：
* `_ACTION_2_ID`：从动作字符串到 ID
* `_ID_2_ACTION`：从 ID 映射回真实牌面字符串
---
## 7. 游戏结束与结算
```python
payoffs = env.get_payoffs()
```
### ✔ 源码对应
| 模块                                | 路径              | 功能         |
| --------------------------------- | --------------- | ---------- |
| `rlcard/games/doudizhu/judger.py` | `judge_payoffs` | 农民对地主的结算方式 |

例如：
* 地主赢 → 地主正分、农民负分
* 农民赢 → 地主负分、农民正分
---
# 📌 示例输出（部分）
```
✨ 斗地主环境初始化完成
[初始手牌分配]
Player 1 [地主]: 4♠ 5♦ 8♥ ...

Player 1 [地主] 动作ID 78 → ['9♠']
Player 2 [农民] 动作ID 0 → ['PASS']
...
🏁 游戏结束
最终得分: [-1, 2, -1]
地主 (Player 1) 得分: 2
🏆 地主获胜!
```
---
# 🧱 可扩展方向建议
本示例可扩展为：
### ✔ 强化学习训练项目
* NFSP
* CFR
### ✔ 对局日志系统
* 保存每一步出牌
* 可视化对局（HTML / GUI）
### ✔ 替换智能体
使用规则智能体 / 学习型智能体：
```
rlcard.agents.dqn_agent
rlcard.agents.nfsp_agent
```
---
# 📜 许可证
本项目依赖 **RLCard (MIT License)**。
你的示例代码也可放置 MIT License 或保持开源。
---