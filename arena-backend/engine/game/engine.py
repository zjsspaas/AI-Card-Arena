import rlcard
from rlcard.utils import set_seed
from typing import Dict, List, Tuple, Optional, Any
import logging
import uuid

logger = logging.getLogger(__name__)

class GameEngine:
    def __init__(self, seed: Optional[int] = None):
        self.seed = seed
        self.env = rlcard.make('doudizhu', config={'seed': self.seed})
        self.game_id = str(uuid.uuid4())
        self._initialized = False
        self._action_id_to_str: Dict[int, str] = {}
        self._action_str_to_id: Dict[str, int] = {}
    
    def reset(self, config: Optional[Dict] = None) -> Tuple[Dict, int]:
        '''
            初始化/重置游戏环境
            
            Args:
                config: RLCard 配置字典，可包含 'seed' 等参数
                
            Returns:
                (初始状态, 当前玩家ID)
            '''
        if self._initialized:
            self.env.reset()
            return self.state, self.player_id
        if self.seed is not None:
            set_seed(self.seed)
        elif config and 'seed' in config:
            set_seed(config['seed'])
        if config is None:
            config = {}
        self.env = rlcard.make('doudizhu', config=config)
        self._init_action_maps()
        
        # 重置环境（发牌、确定地主等）
        state, player_id = self.env.reset()
        
        self._initialized = True
        logger.info(f"游戏环境初始化完成，当前玩家: {player_id}")
        return state, player_id
    def _init_action_maps(self):
        """初始化动作映射表（ID <-> 字符串）"""
        if self.env is None:
            return
        
        # RLCard 的 DoudizhuEnv 内部有 _ID_2_ACTION 和 _ACTION_2_ID 字典
        if hasattr(self.env, '_ID_2_ACTION'):
            self._action_id_to_str = self.env._ID_2_ACTION.copy()
            # 创建反向映射
            self._action_str_to_id = {
                v: k for k, v in self._action_id_to_str.items()
            }
        else:
            logger.warning("无法获取动作映射表，将使用备用方法")
    def get_state(self, player_id: int) -> Dict:
        """
        获取指定玩家的游戏状态
        
        Args:
            player_id: 玩家ID (0, 1, 2)
            
        Returns:
            玩家状态字典，包含：
            - raw_obs: 原始观察数据
            - legal_actions: 合法动作ID列表
            - hand: 手牌
            - trace: 出牌历史
            - others_hand: 其他玩家剩余牌数
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化，请先调用 reset()")
        
        if player_id < 0 or player_id >= self.num_players:
            raise ValueError(f"无效的玩家ID: {player_id}")
        
        return self.env.get_state(player_id)
    def step(self, action: Any) -> Tuple[Dict, int]:
        """
        执行一步动作
        
        Args:
            action: 动作，可以是：
                - int: RLCard 动作ID
                - str: 动作字符串（如 "pass", "33", "AAA"）
                
        Returns:
            (下一步状态, 下一步玩家ID)
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化，请先调用 reset()")
        
        # 如果传入的是字符串，转换为动作ID
        if isinstance(action, str):
            action_id = self.action_str_to_id(action)
            if action_id is None:
                raise ValueError(f"无效的动作字符串: {action}")
        else:
            action_id = action
        
        # 执行动作
        try:
            next_state, next_player_id = self.env.step(action_id)
            return next_state, next_player_id
        except Exception as e:
            logger.error(f"执行动作失败: {action} -> {action_id}, 错误: {str(e)}")
            raise ValueError(f"非法动作: {str(e)}")
    
    def get_legal_actions(self, player_id: int) -> List[str]:
        """
        获取指定玩家的合法动作列表（返回动作字符串）
        
        Args:
            player_id: 玩家ID
            
        Returns:
            合法动作字符串列表，如 ["pass", "33", "44", "AAA"]
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化")
        
        state = self.get_state(player_id)
        legal_action_ids = state.get('legal_actions', [])
        
        # 转换为动作字符串
        legal_actions = []
        for action_id in legal_action_ids:
            action_str = self.action_id_to_str(action_id)
            if action_str:
                legal_actions.append(action_str)
        
        return legal_actions
    def action_id_to_str(self, action_id: int) -> Optional[str]:
        """
        将动作ID转换为动作字符串
        
        Args:
            action_id: RLCard 动作ID
            
        Returns:
            动作字符串，如 "pass", "33", "AAA"
        """
        if self._action_id_to_str:
            return self._action_id_to_str.get(action_id)
        
        # 备用方法：使用 RLCard 的 _decode_action 方法
        if self.env and hasattr(self.env, '_decode_action'):
            try:
                return self.env._decode_action(action_id)
            except:
                pass
        
        return None
    
    def action_str_to_id(self, action_str: str) -> Optional[int]:
        """
        将动作字符串转换为动作ID
        
        Args:
            action_str: 动作字符串，如 "pass", "33", "AAA"
            
        Returns:
            RLCard 动作ID
        """
        # 优先使用映射表
        if self._action_str_to_id:
            return self._action_str_to_id.get(action_str)
        
        # 备用方法：遍历查找
        if self.env and hasattr(self.env, '_ACTION_2_ID'):
            return self.env._ACTION_2_ID.get(action_str)
        
        return None
    def get_landlord_id(self) -> int:
        """
        获取地主玩家ID
        
        Returns:
            地主玩家ID (0, 1, 或 2)
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化")
        
        return self.env.game.round.landlord_id
    def get_perfect_information(self) -> Dict:
        """
        获取完整信息（用于日志、调试等）
        
        Returns:
            包含所有玩家手牌等完整信息的字典
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化")
        
        return self.env.get_perfect_information()
    def is_over(self) -> bool:
        """
        检查游戏是否结束
        
        Returns:
            True 如果游戏已结束，False 否则
        """
        if not self._initialized:
            return False
        
        return self.env.is_over()
    def get_payoffs(self) -> List[float]:
        """
        获取最终得分（游戏结束后调用）
        
        Returns:
            三个玩家的得分列表 [player0_score, player1_score, player2_score]
            地主获胜时地主得正分，农民得负分
            农民获胜时地主得负分，农民得正分
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化")
        
        if not self.is_over():
            logger.warning("游戏未结束，返回当前得分可能不准确")
        
        payoffs = self.env.get_payoffs()
        
        # 转换为列表（RLCard 可能返回 numpy array）
        if hasattr(payoffs, 'tolist'):
            return payoffs.tolist()
        return list(payoffs)
    def get_current_player_id(self) -> int:
        """
        获取当前应该行动的玩家ID
        
        Returns:
            当前玩家ID
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化")
        
        return self.env.get_player_id()
    @property
    def num_players(self) -> int:
        """获取玩家数量"""
        if self.env:
            return self.env.num_players
        return 3  # 斗地主固定3个玩家
    
    @property
    def num_actions(self) -> int:
        """获取动作空间大小"""
        if self.env:
            return self.env.num_actions
        return 0
    
    def get_hand_cards(self, player_id: int) -> List[str]:
        """
        获取指定玩家的手牌（带花色）
        
        Args:
            player_id: 玩家ID
            
        Returns:
            手牌列表，如 ["3♠", "4♥", "5♦"]
        """
        if not self._initialized:
            raise ValueError("游戏环境未初始化")
        
        perfect_info = self.get_perfect_information()
        hands = perfect_info.get('hand_cards_with_suit', [])
        
        if player_id < len(hands):
            # 返回带花色的手牌字符串
            return hands[player_id].split() if isinstance(hands[player_id], str) else hands[player_id]
        
        return []
    
    def get_hand_cards_count(self, player_id: int) -> int:
        """
        获取指定玩家的手牌数量
        
        Args:
            player_id: 玩家ID
            
        Returns:
            手牌数量
        """
        hand_cards = self.get_hand_cards(player_id)
        return len(hand_cards)
    
    def get_all_players_card_count(self) -> List[int]:
        """
        获取所有玩家的手牌数量
        
        Returns:
            [player0_count, player1_count, player2_count]
        """
        return [self.get_hand_cards_count(i) for i in range(self.num_players)]
    
    def cleanup(self):
        """清理资源"""
        self.env = None
        self._initialized = False
        self._action_id_to_str.clear()
        self._action_str_to_id.clear()
        logger.info("游戏引擎资源已清理")