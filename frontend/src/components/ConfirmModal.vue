<template>
  <div v-if="isOpen" class="modal-overlay" @click.self="cancel">
    <div class="modal">
      <div class="modal-header">
        <h3>{{ title }}</h3>
      </div>
      <div class="modal-body">
        <p>{{ message }}</p>
      </div>
      <div class="modal-footer">
        <button @click="cancel" class="btn-secondary">{{ cancelText }}</button>
        <button @click="confirm" :class="danger ? 'btn-danger' : 'btn-primary'">{{ confirmText }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  isOpen: Boolean,
  title: { type: String, default: '确认操作' },
  message: { type: String, required: true },
  confirmText: { type: String, default: '确定' },
  cancelText: { type: String, default: '取消' },
  danger: { type: Boolean, default: false }
})

const emit = defineEmits(['confirm', 'cancel'])

function confirm() {
  emit('confirm')
}

function cancel() {
  emit('cancel')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0; 
  left: 0; 
  right: 0; 
  bottom: 0; 
  background: rgba(0,0,0,0.5); 
  display: flex; 
  align-items: center; 
  justify-content: center; 
  z-index: 2000;
  animation: fadeIn 0.2s;
}

.modal {
  background: white;
  padding: 24px;
  border-radius: 12px;
  width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0,0,0,0.1);
  animation: slideIn 0.2s;
}

.modal-header h3 {
  margin: 0 0 16px;
  font-size: 1.25rem;
  color: #1e293b;
}

.modal-body p {
  color: #64748b;
  font-size: 1rem;
  line-height: 1.5;
  margin-bottom: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-secondary {
  background: white;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
.btn-secondary:hover {
  background: #f1f5f9;
}

.btn-primary {
  background: #2563eb;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
.btn-primary:hover {
  background: #1d4ed8;
}

.btn-danger {
  background: #dc2626;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
}
.btn-danger:hover {
  background: #b91c1c;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { transform: translateY(-20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}
</style>
