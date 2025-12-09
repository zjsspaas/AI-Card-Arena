# Agent输出合约
## 示例
```json
{
  "response_id": "uuid-a8b2c4d6-11e5-4f73-9a00-111222334455",
  "action": "3336",
  "thought": {
    "reasoning": "地主剩3张牌，下家剩10张。我手里对2最大，但为了保留对局势的控制权，先用对K试探地主的牌型，避免过早暴露底牌。",
    "confidence_score": 0.85 
  }
}
```
* response_id：Agent侧生成的唯一响应 ID
* action：当前执行的动作，要求在input-std的actions中字符串选取一个
* reasoning：Agent思考过程
* confidence_score：决策的自信度，范围 0.0 - 1.0