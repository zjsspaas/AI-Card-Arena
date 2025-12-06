import rlcard
from rlcard.agents import RandomAgent
from rlcard.utils import set_seed

# --- 1. 环境设置 ---
# 使用 make 创建环境
env = rlcard.make('doudizhu', config={'seed': 42})
set_seed(42)

# --- 2. 智能体设置 ---
agents = []
# Env类源码中: self.num_players = self.game.get_num_players()
for i in range(env.num_players):
    # Env类源码中: self.num_actions = self.game.get_num_actions()
    agent = RandomAgent(num_actions=env.num_actions)
    agents.append(agent)

env.set_agents(agents)

print("✨ 斗地主环境初始化完成")

# --- 3. 游戏初始化 ---
print("\n--- 🃏 游戏开始 ---")

# Env.reset() 源码返回: (self._extract_state(state), player_id)
state, current_player_id = env.reset()

# --- 4. 获取并打印初始手牌 (基于 provided source code) ---
print("\n[初始手牌分配]")
# 使用 DoudizhuEnv 源码中定义的 get_perfect_information()
perfect_info = env.get_perfect_information()
hands = perfect_info['hand_cards_with_suit']  # 源码中定义了此键

# 获取地主ID (基于 provided source: get_payoffs 使用了 env.game.round.landlord_id)
landlord_id = env.game.round.landlord_id

for i, hand_str in enumerate(hands):
    role = "地主" if i == landlord_id else "农民"
    print(f"Player {i} [{role}]: {hand_str}")

print("\n--- ⚔️ 出牌阶段 ---")

# --- 5. 游戏循环 ---
# Env.is_over() 源码调用 self.game.is_over()
while not env.is_over():
    # 获取当前行动玩家
    current_player_id = env.get_player_id()

    # 获取智能体需要的状态
    state = env.get_state(current_player_id)

    # 智能体决策
    action = agents[current_player_id].step(state)

    # --- 动作解码 (基于 provided source code) ---
    # DoudizhuEnv.__init__ 中定义了 self._ID_2_ACTION
    # DoudizhuEnv._decode_action 返回 self._ID_2_ACTION[action_id]
    # 我们直接访问字典，或者调用内部方法
    if hasattr(env, '_ID_2_ACTION'):
        action_card_str = env._ID_2_ACTION[action]
    else:
        # 备用：如果直接访问受限，尝试调用 _decode_action
        action_card_str = env._decode_action(action)

    # 美化打印
    role_label = "地主" if current_player_id == landlord_id else "农民"
    print(f"Player {current_player_id} [{role_label}] 动作ID {action} ➡️ {action_card_str}")

    # 环境推进一步
    # Env.step() 返回 (next_state, next_player_id)
    next_state, next_player_id = env.step(action)

# --- 6. 游戏结算 ---
print("\n--- 🏁 游戏结束 ---")
# DoudizhuEnv.get_payoffs() 源码调用 self.game.judger.judge_payoffs
payoffs = env.get_payoffs()

print(f"最终得分: {payoffs}")
print(f"地主 (Player {landlord_id}) 得分: {payoffs[landlord_id]}")

if payoffs[landlord_id] > 0:
    print("🏆 结果: 地主获胜!")
else:
    print("🏆 结果: 农民获胜!")