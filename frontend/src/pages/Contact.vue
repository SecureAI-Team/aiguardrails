<template>
  <LandingLayout>
    <div class="page-container">
      <div class="header">
        <h1>联系我们</h1>
        <p>有任何问题或建议？我们随时倾听。</p>
      </div>

      <div class="contact-grid">
        <div class="contact-info">
          <div class="info-item">
            <h3>📍 地址</h3>
            <p>北京市海淀区中关村软件园 1 号楼</p>
          </div>
          <div class="info-item">
            <h3>📧 邮箱</h3>
            <p>support@aiguardrails.com</p>
            <p>sales@aiguardrails.com</p>
          </div>
          <div class="info-item">
            <h3>📞 电话</h3>
            <p>+86 10-1234-5678</p>
          </div>
        </div>

        <div class="contact-form-card">
          <form @submit.prevent="submitForm">
            <div class="form-group">
              <label>姓名</label>
              <input v-model="form.name" type="text" placeholder="您的姓名" required />
            </div>
            <div class="form-group">
              <label>邮箱</label>
              <input v-model="form.email" type="email" placeholder="your@email.com" required />
            </div>
            <div class="form-group">
              <label>主题</label>
              <input v-model="form.subject" type="text" placeholder="咨询主题" required />
            </div>
            <div class="form-group">
              <label>留言内容</label>
              <textarea v-model="form.message" rows="5" placeholder="请详细描述您的需求..." required></textarea>
            </div>
            <button type="submit" class="btn-submit" :disabled="sending">
              {{ sending ? '发送中...' : '发送留言' }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </LandingLayout>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import LandingLayout from '../components/LandingLayout.vue'

const form = reactive({
  name: '',
  email: '',
  subject: '',
  message: ''
})
const sending = ref(false)

function submitForm() {
  sending.value = true
  setTimeout(() => {
    alert('感谢您的留言！我们会尽快与您联系。')
    sending.value = false
    form.name = ''
    form.email = ''
    form.subject = ''
    form.message = ''
  }, 1000)
}
</script>

<style scoped>
.page-container {
  padding: 60px 20px;
  max-width: 1000px;
  margin: 0 auto;
}
.header {
  text-align: center;
  margin-bottom: 60px;
}
.header h1 { font-size: 2.5rem; color: #1e293b; margin-bottom: 16px; }
.header p { color: #64748b; font-size: 1.2rem; }

.contact-grid {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 60px;
}

.contact-info .info-item { margin-bottom: 40px; }
.contact-info h3 { font-size: 1.2rem; color: #1e293b; margin-bottom: 12px; }
.contact-info p { color: #64748b; margin: 4px 0; font-size: 1.1rem; }

.contact-form-card {
  background: white;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
  border: 1px solid #e2e8f0;
}

.form-group { margin-bottom: 24px; }
.form-group label { display: block; margin-bottom: 8px; font-weight: 500; color: #475569; }
.form-group input, .form-group textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.2s;
}
.form-group input:focus, .form-group textarea:focus {
  outline: none;
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.btn-submit {
  width: 100%;
  padding: 14px;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-submit:hover:not(:disabled) {
  background: #1d4ed8;
}
.btn-submit:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .contact-grid { grid-template-columns: 1fr; gap: 40px; }
}
</style>
