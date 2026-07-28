/* ============================================
   NATS Framework - Main JavaScript
   ============================================ */

// ============================================
// DOM Ready
// ============================================
document.addEventListener('DOMContentLoaded', function() {
    initApp();
});

function initApp() {
    initHTMX();
    initAlpine();
    initNotifications();
    initModals();
    initForms();
}

// ============================================
// HTMX Configuration
// ============================================
function initHTMX() {
    if (typeof htmx !== 'undefined') {
        htmx.config.historyCacheSize = 20;
        htmx.config.defaultSwapStyle = 'innerHTML';
        htmx.config.defaultSwapDelay = 0;
        htmx.config.defaultTimeout = 10000;
        
        // إضافة حدث للطلبات الناجحة
        document.body.addEventListener('htmx:afterRequest', function(evt) {
            if (evt.detail.successful) {
                // معالجة الطلبات الناجحة
                showNotification('Success', 'Operation completed successfully', 'success');
            }
        });
        
        // إضافة حدث للأخطاء
        document.body.addEventListener('htmx:responseError', function(evt) {
            let errorMsg = 'An error occurred';
            if (evt.detail.xhr && evt.detail.xhr.responseText) {
                try {
                    const response = JSON.parse(evt.detail.xhr.responseText);
                    if (response.error) {
                        errorMsg = response.error;
                    }
                } catch(e) {}
            }
            showNotification('Error', errorMsg, 'danger');
        });
        
        // إضافة حدث للتحقق من الصلاحية
        document.body.addEventListener('htmx:validation:validate', function(evt) {
            const form = evt.detail.elt;
            const inputs = form.querySelectorAll('[required]');
            let isValid = true;
            
            inputs.forEach(function(input) {
                if (!input.value.trim()) {
                    input.classList.add('is-invalid');
                    isValid = false;
                } else {
                    input.classList.remove('is-invalid');
                }
            });
            
            return isValid;
        });
    }
}

// ============================================
// Alpine.js Configuration
// ============================================
function initAlpine() {
    if (typeof Alpine !== 'undefined') {
        Alpine.data('app', () => ({
            theme: localStorage.getItem('theme') || 'light',
            sidebarOpen: true,
            notifications: [],
            
            init() {
                this.applyTheme();
                this.loadNotifications();
            },
            
            toggleTheme() {
                this.theme = this.theme === 'light' ? 'dark' : 'light';
                localStorage.setItem('theme', this.theme);
                this.applyTheme();
            },
            
            applyTheme() {
                document.documentElement.classList.toggle('dark', this.theme === 'dark');
            },
            
            toggleSidebar() {
                this.sidebarOpen = !this.sidebarOpen;
            },
            
            loadNotifications() {
                fetch('/api/notifications')
                    .then(response => response.json())
                    .then(data => {
                        if (data.success) {
                            this.notifications = data.data || [];
                        }
                    })
                    .catch(() => {});
            },
            
            addNotification(message, type = 'info') {
                this.notifications.unshift({
                    id: Date.now(),
                    message: message,
                    type: type,
                    timestamp: new Date().toISOString()
                });
            },
            
            removeNotification(id) {
                this.notifications = this.notifications.filter(n => n.id !== id);
            }
        }));
    }
}

// ============================================
// Notifications
// ============================================
function initNotifications() {
    // إضافة حدث للإشعارات من الخادم
    if (typeof EventSource !== 'undefined') {
        const eventSource = new EventSource('/api/notifications/stream');
        eventSource.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                showNotification(data.title || 'Notification', data.message || '', data.type || 'info');
            } catch(e) {}
        };
    }
}

function showNotification(title, message, type = 'info') {
    const container = document.getElementById('notifications');
    if (!container) return;
    
    const colors = {
        success: 'bg-green-100 border-green-400 text-green-700',
        danger: 'bg-red-100 border-red-400 text-red-700',
        warning: 'bg-yellow-100 border-yellow-400 text-yellow-700',
        info: 'bg-blue-100 border-blue-400 text-blue-700'
    };
    
    const icons = {
        success: 'fa-check-circle',
        danger: 'fa-exclamation-circle',
        warning: 'fa-exclamation-triangle',
        info: 'fa-info-circle'
    };
    
    const notification = document.createElement('div');
    notification.className = `border-l-4 p-4 mb-2 rounded shadow ${colors[type] || colors.info}`;
    notification.innerHTML = `
        <div class="flex items-start">
            <div class="flex-shrink-0">
                <i class="fas ${icons[type] || icons.info} mt-1"></i>
            </div>
            <div class="ml-3 flex-1">
                <p class="font-bold">${title}</p>
                <p class="text-sm">${message}</p>
            </div>
            <button onclick="this.parentElement.parentElement.remove()" class="flex-shrink-0 ml-4">
                <i class="fas fa-times"></i>
            </button>
        </div>
    `;
    
    container.prepend(notification);
    
    // إزالة الإشعار تلقائياً بعد 5 ثواني
    setTimeout(() => {
        notification.remove();
    }, 5000);
}

// ============================================
// Modals
// ============================================
function initModals() {
    document.querySelectorAll('[data-modal-toggle]').forEach(button => {
        button.addEventListener('click', function() {
            const target = this.getAttribute('data-modal-target');
            const modal = document.getElementById(target);
            if (modal) {
                modal.classList.toggle('hidden');
            }
        });
    });
    
    // إغلاق المودال عند النقر خارج المحتوى
    document.querySelectorAll('.modal-overlay').forEach(overlay => {
        overlay.addEventListener('click', function(e) {
            if (e.target === this) {
                this.classList.add('hidden');
            }
        });
    });
}

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.remove('hidden');
        document.body.style.overflow = 'hidden';
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('hidden');
        document.body.style.overflow = 'auto';
    }
}

// ============================================
// Forms
// ============================================
function initForms() {
    // التحقق من النماذج
    document.querySelectorAll('form[data-validate]').forEach(form => {
        form.addEventListener('submit', function(e) {
            const inputs = this.querySelectorAll('[required]');
            let isValid = true;
            
            inputs.forEach(function(input) {
                if (!input.value.trim()) {
                    input.classList.add('is-invalid');
                    isValid = false;
                } else {
                    input.classList.remove('is-invalid');
                }
            });
            
            if (!isValid) {
                e.preventDefault();
                showNotification('Validation Error', 'Please fill in all required fields', 'danger');
            }
        });
    });
    
    // تحقق فوري من الحقول
    document.querySelectorAll('[required]').forEach(input => {
        input.addEventListener('blur', function() {
            if (this.value.trim()) {
                this.classList.remove('is-invalid');
                this.classList.add('is-valid');
            } else {
                this.classList.remove('is-valid');
                this.classList.add('is-invalid');
            }
        });
    });
}

// ============================================
// Utilities
// ============================================
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function formatNumber(number) {
    return new Intl.NumberFormat('en-US').format(number);
}

function formatCurrency(amount, currency = 'USD') {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: currency
    }).format(amount);
}

function truncate(text, length = 50) {
    if (text.length <= length) return text;
    return text.substring(0, length) + '...';
}

function debounce(func, wait = 300) {
    let timeout;
    return function(...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    };
}

// ============================================
// Export utilities (if module)
// ============================================
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        showNotification,
        openModal,
        closeModal,
        formatDate,
        formatNumber,
        formatCurrency,
        truncate,
        debounce
    };
}