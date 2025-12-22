<template>
  <div class="replay-collection">
    <!-- 主标题区域 - 与回放,首页页面一致 -->
    <div class="page-header">
      <h1>斗地主对战精彩回放集锦</h1>
      <p class="subtitle">精选高光时刻，领略斗地主绝妙策略</p>
    </div>

    <!-- 回放视频网格布局 -->
    <div class="replay-grid">
      <div 
        v-for="(video, index) in featuredVideos" 
        :key="index" 
        class="video-item"
        @click="goToVideo(video.id)"
      >
        <div class="video-thumbnail">
          <img :src="video.thumbnail" :alt="video.title">
          <div class="video-duration">{{ formatDuration(video.duration) }}</div>
          <div class="video-hover">
            <div class="play-icon">▶</div>
            <div class="video-info-hover">
              <p class="video-title">{{ video.title }}</p>
              <div class="video-stats">
                <span class="views">观看: {{ video.views }}万</span>
                <span class="likes">点赞: {{ video.likes }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载更多按钮 -->
    <div class="load-more">
      <button class="btn btn-primary">加载更多精彩回放</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// 回放视频数据
const featuredVideos = ref([
  {
    id: 1,
    title: '神级预判！地主明牌AA仍惨败，农民教科书级防守',
    thumbnail: '/images/suoluetu.jpg',
    duration: 235, // 秒
    views: 158,
    likes: 423
  },
  {
    id: 2,
    title: '极限反杀！农民双王炸绝地反击，全场观众沸腾',
    thumbnail: '/images/suoluetu.jpg',
    duration: 197,
    views: 284,
    likes: 612
  },
  {
    id: 3,
    title: '大师对决：千术解密之精准叫地主技巧',
    thumbnail: '/images/suoluetu.jpg',
    duration: 412,
    views: 327,
    likes: 821
  },
  {
    id: 4,
    title: '欢乐时刻：朋友聚会搞笑斗地主，笑翻全场',
    thumbnail: '/images/suoluetu.jpg',
    duration: 164,
    views: 415,
    likes: 953
  },
  {
    id: 5,
    title: '心理战：记牌器高手的完美控分战术解析',
    thumbnail: '/images/suoluetu.jpg',
    duration: 301,
    views: 276,
    likes: 598
  },
  {
    id: 6,
    title: '新手必学：从零开始的斗地主进阶之路',
    thumbnail: 'https://via.placeholder.com/300x180?text=斗地主教学2',
    duration: 523,
    views: 642,
    likes: 1284
  }
])

// 格式化时长
const formatDuration = (seconds) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// 跳转到视频详情页- 使用路由跳转
const goToVideo = (id) => {
  router.push({ 
    name: 'ReplayDoudizhu_detail',  // 必须使用路由名称
    params: { videoId: id }         // 传递参数
  })
}
</script>

<style scoped>
.replay-collection {
  padding: 20px 20px 60px;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #f5f7fa;
}

/* 主标题样式 - 与示例页面一致 */
.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-header h1 {
  font-size: 36px;
  font-weight: 700;
  color: #333;
  margin-bottom: 10px;
  letter-spacing: 1px;
}

.subtitle {
  color: #666;
  font-size: 18px;
  font-weight: 400;
  margin: 0;
}

/* 回放视频网格 */
.replay-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  width: 100%;
}

.video-item {
  background-color: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.video-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 25px rgba(0, 0, 0, 0.15);
}

.video-thumbnail {
  position: relative;
  padding-top: 56.25%; /* 16:9比例 */
  overflow: hidden;
}

.video-thumbnail img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.video-item:hover .video-thumbnail img {
  transform: scale(1.05);
}

.video-duration {
  position: absolute;
  bottom: 10px;
  right: 10px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.video-hover {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.video-item:hover .video-hover {
  opacity: 1;
}

.play-icon {
  font-size: 48px;
  color: white;
  margin-bottom: 15px;
}

.video-info-hover {
  text-align: center;
}

.video-title {
  font-size: 25px;
  font-weight: 600;
  color: white;
  margin-bottom: 8px;
  line-height: 1.4;
}

.video-stats {
  display: flex;
  justify-content: center;
  gap: 15px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
}

.load-more {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

.btn-primary {
  padding: 12px 30px;
  background-color: #1e3a8a;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.btn-primary:hover {
  background-color: #152c6e;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header h1 {
    font-size: 28px;
  }
  
  .replay-grid {
    grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    gap: 20px;
  }
  
  .video-title {
    font-size: 14px;
    margin-bottom: 5px;
  }
  
  .video-stats {
    font-size: 12px;
    gap: 8px;
  }
}
</style>