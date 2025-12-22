<template>
  <div class="video-detail">
    <!-- 返回按钮 -->
    <div class="back-button">
      <button @click="goBack" class="btn-back">← 返回列表</button>
    </div>

    <!-- 视频播放区域 -->
    <div class="video-player-section">
      <div class="video-container">
        <video 
          :src="videoUrl" 
          controls 
          width="100%" 
          height="500"
          poster="./images/suoluetu.jpg"
        >
          您的浏览器不支持视频播放
        </video>
      </div>
      
      <!-- 视频信息 -->
      <div class="video-info">
        <h1>{{ videoData?.title }}</h1>
        <div class="video-meta">
          <span class="views">观看: {{ videoData?.views }}万次</span>
          <span class="likes">点赞: {{ videoData?.likes }}</span>
          <span class="duration">时长: {{ formatDuration(videoData?.duration) }}</span>
          <span class="upload-time">上传时间: {{ videoData?.uploadTime }}</span>
        </div>
        <div class="video-description">
          <h3>视频简介</h3>
          <p>{{ videoData?.description }}</p>
        </div>
      </div>
    </div>

    <!-- 相关推荐 -->
    <div class="related-videos">
      <h2>相关推荐</h2>
      <div class="related-grid">
        <div 
          v-for="related in relatedVideos" 
          :key="related.id"
          class="related-item"
          @click="goToVideo(related.id)"
        >
          <img :src="related.thumbnail" :alt="related.title">
          <div class="related-info">
            <h4>{{ related.title }}</h4>
            <span class="related-views">{{ related.views }}万观看</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

// 接收路由参数
const videoId = ref(route.params.videoId)

// 视频数据
const videoData = ref(null)
const videoUrl = ref('') // 实际项目中这里应该是视频文件的URL

// 相关视频数据
const relatedVideos = ref([
  {
    id: 2,
    title: '极限反杀！农民双王炸绝地反击',
    thumbnail: './images/suoluetu.jpg',
    views: 284
  },
  // ... 其他相关视频
])

// 模拟根据ID获取视频数据
const fetchVideoData = async (id) => {
  // 这里应该是实际的API调用
  // 现在使用模拟数据
  const mockData = {
    1: {
      id: 1,
      title: '神级预判！地主明牌AA仍惨败，农民教科书级防守',
      thumbnail: '/images/suoluetu.image',
      duration: 235,
      views: 158,
      likes: 423,
      uploadTime: '2小时前',
      description: '本局比赛中，农民玩家展现了惊人的预判能力，在地主明牌AA的情况下，通过精妙的配合和准确的记牌，成功实现逆转。这局比赛堪称教科书级别的防守范例...'
    },
    2: {
      id: 2,
      title: '极限反杀！农民双王炸绝地反击，全场观众沸腾',
      thumbnail: './images/suoluetu.jpg',
      duration: 197,
      views: 284,
      likes: 612,
      uploadTime: '5小时前',
      description: '在看似必败的局面下，农民玩家利用双王炸实现了惊天逆转，展现了斗地主游戏的无限可能...'
    }
  }
  
  videoData.value = mockData[id] || null
  // 设置视频URL（实际项目中应该从后端获取）
  videoUrl.value = `/videos/doudi_zhu_${id}.mp4` // 假设的视频文件路径
}

// 返回上一页
const goBack = () => {
  router.go(-1) // 返回上一页
  // 或者 router.push('/replay/doudizhu')
}

// 跳转到其他视频
const goToVideo = (id) => {
  router.push({ 
    name: 'VideoDetail', 
    params: { videoId: id } 
  })
}

// 格式化时长
const formatDuration = (seconds) => {
  if (!seconds) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// 组件挂载时获取数据
onMounted(() => {
  fetchVideoData(videoId.value)
})
</script>

<style scoped>
.video-detail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.back-button {
  margin-bottom: 20px;
}

.btn-back {
  background: #6c757d;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.video-player-section {
  display: flex;
  flex-direction: column;
  gap: 30px;
  margin-bottom: 50px;
}

.video-container {
  width: 100%;
  background: #000;
  border-radius: 8px;
  overflow: hidden;
}

.video-info h1 {
  font-size: 28px;
  margin-bottom: 15px;
  color: #333;
}

.video-meta {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
  color: #666;
  font-size: 14px;
}

.video-description {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
}

.video-description h3 {
  margin-bottom: 10px;
  color: #333;
}

.related-videos h2 {
  margin-bottom: 20px;
  color: #333;
}

.related-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
}

.related-item {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: transform 0.3s;
}

.related-item:hover {
  transform: translateY(-3px);
}

.related-item img {
  width: 100%;
  height: 120px;
  object-fit: cover;
}

.related-info {
  padding: 10px;
}

.related-info h4 {
  font-size: 14px;
  margin-bottom: 5px;
  color: #333;
}

.related-views {
  font-size: 12px;
  color: #666;
}
</style>