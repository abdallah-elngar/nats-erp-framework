<!DOCTYPE html>
<html lang="{{.Locale}}" x-data="{ theme: localStorage.getItem('theme') || 'light' }" :class="theme">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} | NATS ERP</title>
    
    <!-- SEO Meta Tags -->
    <meta name="description" content="{{.Description}}">
    <meta name="keywords" content="ERP, NATS, Framework, Go">
    <meta name="author" content="NATS Framework">
    
    <!-- CSRF Token -->
    <meta name="csrf-token" content="{{.CSRFToken}}">
    
    <!-- Static Files - Local (No CDN) -->
    <link rel="stylesheet" href="/static/css/tailwind/tailwind.min.css">
    <link rel="stylesheet" href="/static/fonts/fontawesome/css/all.min.css">
    <link rel="stylesheet" href="/static/css/app.css">
    <link rel="stylesheet" href="/static/css/admin.css">
    <link rel="stylesheet" href="/static/css/themes/{{.Theme}}.css">
    
    <!-- Favicon -->
    <link rel="icon" type="image/x-icon" href="/static/images/favicon.ico">
    <link rel="apple-touch-icon" href="/static/images/logo.svg">
    
    <!-- HTMX - Dynamic UI Updates -->
    <script src="/static/js/htmx/htmx.min.js"></script>
    
    <!-- Alpine.js - Reactive Components -->
    <script src="/static/js/alpine/alpine.min.js"></script>
    
    <!-- Application JavaScript -->
    <script src="/static/js/app.js"></script>
    <script src="/static/js/admin.js"></script>
    
    <!-- HTMX Configuration -->
    <script>
        document.addEventListener('DOMContentLoaded', function() {
            if (typeof htmx !== 'undefined') {
                htmx.config.historyCacheSize = 50;
                htmx.config.defaultSwapStyle = 'innerHTML';
                htmx.config.defaultSwapDelay = 0;
                htmx.config.defaultTimeout = 30000;
                htmx.config.selfRequestsOnly = false;
                htmx.config.allowEval = true;
                htmx.config.useTemplateFragments = true;
                
                // إضافة التوكن CSRF إلى جميع الطلبات
                const csrfToken = document.querySelector('meta[name="csrf-token"]')?.getAttribute('content');
                if (csrfToken) {
                    document.body.addEventListener('htmx:configRequest', function(evt) {
                        evt.detail.headers['X-CSRF-Token'] = csrfToken;
                    });
                }
            }
        });
    </script>
    
    <!-- Custom Styles -->
    <style>
        /* إعدادات إضافية مخصصة */
        [x-cloak] { display: none !important; }
        
        .transition-all { transition: all 0.3s ease-in-out; }
        .transition-colors { transition: background-color 0.2s, color 0.2s, border-color 0.2s; }
        
        /* شريط التمرير المخصص */
        ::-webkit-scrollbar { width: 8px; height: 8px; }
        ::-webkit-scrollbar-track { background: #f1f1f1; border-radius: 4px; }
        ::-webkit-scrollbar-thumb { background: #c1c1c1; border-radius: 4px; }
        ::-webkit-scrollbar-thumb:hover { background: #a8a8a8; }
        
        /* وضع الظلام */
        .dark ::-webkit-scrollbar-track { background: #1a1a2e; }
        .dark ::-webkit-scrollbar-thumb { background: #2d2d44; }
        .dark ::-webkit-scrollbar-thumb:hover { background: #3d3d5c; }
        
        /* تحميل HTMX */
        .htmx-indicator { opacity: 0; transition: opacity 500ms ease-in; }
        .htmx-request .htmx-indicator { opacity: 1; }
        .htmx-request.htmx-indicator { opacity: 1; }
        
        /* تأثيرات النقر */
        .click-ripple {
            position: relative;
            overflow: hidden;
        }
        .click-ripple::after {
            content: '';
            position: absolute;
            border-radius: 50%;
            background: rgba(255,255,255,0.3);
            transform: scale(0);
            animation: ripple 0.6s linear;
            pointer-events: none;
        }
        @keyframes ripple {
            to { transform: scale(4); opacity: 0; }
        }
    </style>
</head>

<body>
    <div id="app" 
         x-data="{
             sidebarOpen: JSON.parse(localStorage.getItem('sidebarOpen') || 'true'),
             showModal: false,
             modalContent: '',
             notifications: [],
             user: null,
             
             init() {
                 // تحميل بيانات المستخدم
                 this.loadUser();
                 
                 // مراقبة تغييرات الجلسة
                 this.$watch('sidebarOpen', value => {
                     localStorage.setItem('sidebarOpen', JSON.stringify(value));
                 });
                 
                 // مراقبة الإشعارات من HTMX
                 document.addEventListener('htmx:afterRequest', (evt) => {
                     const trigger = evt.detail.xhr?.getResponseHeader('HX-Trigger');
                     if (trigger) {
                         try {
                             const data = JSON.parse(trigger);
                             if (data.notification) {
                                 this.addNotification(data.notification.message, data.notification.type);
                             }
                         } catch(e) {}
                     }
                 });
             },
             
             loadUser() {
                 fetch('/api/auth/profile')
                     .then(res => res.json())
                     .then(data => {
                         if (data.success) {
                             this.user = data.data;
                         }
                     })
                     .catch(() => {});
             },
             
             toggleSidebar() {
                 this.sidebarOpen = !this.sidebarOpen;
             },
             
             addNotification(message, type = 'info') {
                 const id = Date.now();
                 this.notifications.unshift({ id, message, type, timestamp: new Date() });
                 
                 // إزالة الإشعار تلقائياً بعد 5 ثواني
                 setTimeout(() => {
                     this.removeNotification(id);
                 }, 5000);
             },
             
             removeNotification(id) {
                 this.notifications = this.notifications.filter(n => n.id !== id);
             },
             
             openModal(content) {
                 this.modalContent = content;
                 this.showModal = true;
                 document.body.style.overflow = 'hidden';
             },
             
             closeModal() {
                 this.showModal = false;
                 this.modalContent = '';
                 document.body.style.overflow = 'auto';
             },
             
             logout() {
                 fetch('/api/auth/logout', { method: 'POST' })
                     .then(() => {
                         window.location.href = '/login';
                     })
                     .catch(() => {
                         window.location.href = '/login';
                     });
             }
         }"
         @keydown.escape="closeModal()"
         class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors">
        
        <!-- ============================================ -->
        <!-- NAVBAR -->
        <!-- ============================================ -->
        <nav class="fixed top-0 left-0 right-0 z-50 bg-white dark:bg-gray-800 shadow-sm transition-colors">
            <div class="flex items-center justify-between h-16 px-4">
                <!-- Left Side -->
                <div class="flex items-center gap-3">
                    <!-- Toggle Sidebar -->
                    <button @click="toggleSidebar()" 
                            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                        <i class="fas fa-bars text-gray-700 dark:text-gray-300 text-xl"></i>
                    </button>
                    
                    <!-- Logo -->
                    <a href="/" class="flex items-center gap-2">
                        <img src="/static/images/logo.svg" alt="NATS" class="h-8 w-8">
                        <span class="text-xl font-bold text-primary-600 dark:text-primary-400 hidden sm:block">
                            NATS ERP
                        </span>
                    </a>
                </div>
                
                <!-- Center - Search -->
                <div class="hidden md:flex flex-1 max-w-md mx-4">
                    <div class="relative w-full">
                        <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
                            <i class="fas fa-search"></i>
                        </span>
                        <input type="text" 
                               placeholder="Search..."
                               class="w-full pl-10 pr-4 py-2 bg-gray-100 dark:bg-gray-700 border border-transparent rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:bg-white dark:focus:bg-gray-600 transition-all"
                               hx-get="/api/search"
                               hx-trigger="keyup changed delay:300ms"
                               hx-target="#search-results"
                               hx-indicator="#search-loading">
                        <div id="search-results" class="absolute top-full left-0 right-0 mt-1 bg-white dark:bg-gray-800 rounded-lg shadow-lg hidden"></div>
                    </div>
                </div>
                
                <!-- Right Side -->
                <div class="flex items-center gap-2">
                    <!-- Theme Toggle -->
                    <button @click="theme = theme === 'light' ? 'dark' : 'light'; localStorage.setItem('theme', theme)" 
                            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                        <i class="fas fa-moon text-gray-700 dark:text-gray-300 text-lg" x-show="theme === 'light'"></i>
                        <i class="fas fa-sun text-gray-700 dark:text-gray-300 text-lg" x-show="theme === 'dark'"></i>
                    </button>
                    
                    <!-- Notifications -->
                    <div class="relative" x-data="{ open: false }">
                        <button @click="open = !open" 
                                class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors relative">
                            <i class="fas fa-bell text-gray-700 dark:text-gray-300 text-lg"></i>
                            <span class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full" 
                                  x-show="notifications.length > 0"></span>
                        </button>
                        
                        <!-- Notifications Dropdown -->
                        <div x-show="open" 
                             @click.away="open = false"
                             class="absolute right-0 mt-2 w-80 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 max-h-96 overflow-y-auto">
                            <div class="p-3 border-b border-gray-200 dark:border-gray-700">
                                <h3 class="font-semibold text-gray-900 dark:text-gray-100">Notifications</h3>
                            </div>
                            <div class="divide-y divide-gray-200 dark:divide-gray-700">
                                <template x-for="notification in notifications" :key="notification.id">
                                    <div class="p-3 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
                                        <p class="text-sm text-gray-900 dark:text-gray-100" x-text="notification.message"></p>
                                        <span class="text-xs text-gray-500" x-text="new Date(notification.timestamp).toLocaleString()"></span>
                                    </div>
                                </template>
                                <template x-if="notifications.length === 0">
                                    <div class="p-4 text-center text-gray-500 dark:text-gray-400">
                                        No notifications
                                    </div>
                                </template>
                            </div>
                        </div>
                    </div>
                    
                    <!-- User Menu -->
                    <div class="relative" x-data="{ open: false }">
                        <button @click="open = !open" 
                                class="flex items-center gap-2 p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <div class="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white font-bold">
                                <span x-text="user ? user.full_name?.charAt(0) || user.username?.charAt(0) || 'U' : 'U'"></span>
                            </div>
                            <span class="hidden sm:block text-sm text-gray-700 dark:text-gray-300" x-text="user ? user.full_name || user.username : 'Guest'"></span>
                            <i class="fas fa-chevron-down text-xs text-gray-500"></i>
                        </button>
                        
                        <!-- Dropdown -->
                        <div x-show="open" 
                             @click.away="open = false"
                             class="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700">
                            <div class="py-1">
                                <a href="/profile" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                                    <i class="fas fa-user mr-2"></i>Profile
                                </a>
                                <a href="/settings" class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                                    <i class="fas fa-cog mr-2"></i>Settings
                                </a>
                                <hr class="my-1 border-gray-200 dark:border-gray-700">
                                <button @click="logout()" class="w-full text-left block px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                                    <i class="fas fa-sign-out-alt mr-2"></i>Logout
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
        
        <!-- ============================================ -->
        <!-- SIDEBAR -->
        <!-- ============================================ -->
        <aside x-show="sidebarOpen"
               x-transition:enter="transition ease-out duration-300"
               x-transition:enter-start="-translate-x-full"
               x-transition:enter-end="translate-x-0"
               x-transition:leave="transition ease-in duration-300"
               x-transition:leave-start="translate-x-0"
               x-transition:leave-end="-translate-x-full"
               class="fixed top-16 left-0 bottom-0 w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 overflow-y-auto z-40 transition-colors">
            
            <nav class="p-4">
                <ul class="space-y-1">
                    <!-- Dashboard -->
                    <li>
                        <a href="/" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-home text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Dashboard</span>
                        </a>
                    </li>
                    
                    <!-- Users -->
                    <li>
                        <a href="/users" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-users text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Users</span>
                        </a>
                    </li>
                    
                    <!-- Roles -->
                    <li>
                        <a href="/roles" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-user-tag text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Roles</span>
                        </a>
                    </li>
                    
                    <!-- Permissions -->
                    <li>
                        <a href="/permissions" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-lock text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Permissions</span>
                        </a>
                    </li>
                    
                    <!-- Apps -->
                    <li>
                        <a href="/apps" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-th-large text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Applications</span>
                        </a>
                    </li>
                    
                    <!-- Settings -->
                    <li>
                        <a href="/settings" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-cog text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Settings</span>
                        </a>
                    </li>
                    
                    <hr class="my-4 border-gray-200 dark:border-gray-700">
                    
                    <!-- Reports -->
                    <li>
                        <a href="/reports" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-file-alt text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Reports</span>
                        </a>
                    </li>
                    
                    <!-- Backups -->
                    <li>
                        <a href="/backups" class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
                            <i class="fas fa-database text-gray-500 dark:text-gray-400 w-5"></i>
                            <span class="text-gray-700 dark:text-gray-300">Backups</span>
                        </a>
                    </li>
                </ul>
            </nav>
        </aside>
        
        <!-- ============================================ -->
        <!-- MAIN CONTENT -->
        <!-- ============================================ -->
        <main class="pt-16 min-h-screen transition-all duration-300"
              :class="sidebarOpen ? 'ml-64' : 'ml-0'">
            
            <!-- Notifications Container -->
            <div id="notifications" 
                 class="fixed top-20 right-4 z-50 w-96 space-y-2 pointer-events-none">
                <template x-for="notification in notifications" :key="notification.id">
                    <div class="pointer-events-auto p-4 rounded-lg shadow-lg border transition-all"
                         :class="{
                             'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-800': notification.type === 'success',
                             'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800': notification.type === 'danger',
                             'bg-yellow-50 dark:bg-yellow-900/20 border-yellow-200 dark:border-yellow-800': notification.type === 'warning',
                             'bg-blue-50 dark:bg-blue-900/20 border-blue-200 dark:border-blue-800': notification.type === 'info'
                         }"
                         x-show="true"
                         x-transition:enter="transition ease-out duration-300"
                         x-transition:enter-start="opacity-0 transform translate-x-full"
                         x-transition:enter-end="opacity-100 transform translate-x-0"
                         x-transition:leave="transition ease-in duration-300"
                         x-transition:leave-start="opacity-100 transform translate-x-0"
                         x-transition:leave-end="opacity-0 transform translate-x-full">
                        <div class="flex items-start justify-between">
                            <div class="flex-1">
                                <p class="text-sm font-medium" 
                                   :class="{
                                       'text-green-800 dark:text-green-200': notification.type === 'success',
                                       'text-red-800 dark:text-red-200': notification.type === 'danger',
                                       'text-yellow-800 dark:text-yellow-200': notification.type === 'warning',
                                       'text-blue-800 dark:text-blue-200': notification.type === 'info'
                                   }"
                                   x-text="notification.message"></p>
                            </div>
                            <button @click="removeNotification(notification.id)" 
                                    class="ml-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors">
                                <i class="fas fa-times"></i>
                            </button>
                        </div>
                    </div>
                </template>
            </div>
            
            <!-- Breadcrumb -->
            <div class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-6 py-3 transition-colors">
                <div class="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                    {{template "breadcrumb" .}}
                </div>
            </div>
            
            <!-- Page Content -->
            <div class="p-6">
                {{template "content" .}}
            </div>
            
            <!-- Loading Indicator -->
            <div id="loading-indicator" 
                 class="fixed bottom-4 right-4 bg-white dark:bg-gray-800 rounded-lg shadow-lg p-3 hidden">
                <div class="flex items-center gap-2">
                    <svg class="animate-spin h-5 w-5 text-primary-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span class="text-sm text-gray-700 dark:text-gray-300">Loading...</span>
                </div>
            </div>
        </main>
        
        <!-- ============================================ -->
        <!-- MODAL -->
        <!-- ============================================ -->
        <div x-show="showModal" 
             x-cloak
             class="fixed inset-0 z-50 overflow-y-auto"
             @keydown.escape="closeModal()">
            
            <!-- Backdrop -->
            <div class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm transition-opacity"
                 @click="closeModal()"></div>
            
            <!-- Modal Content -->
            <div class="flex items-center justify-center min-h-screen px-4 py-8">
                <div class="relative bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-y-auto transition-all"
                     @click.away="closeModal()">
                    
                    <!-- Modal Body -->
                    <div id="modal-content" class="p-6">
                        <!-- سيتم تعبئته بواسطة HTMX -->
                    </div>
                </div>
            </div>
        </div>
        
        <!-- ============================================ -->
        <!-- FOOTER -->
        <!-- ============================================ -->
        <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 transition-colors mt-auto">
            <div class="container mx-auto px-6 py-4">
                <div class="flex flex-col sm:flex-row justify-between items-center gap-2">
                    <div class="text-sm text-gray-600 dark:text-gray-400">
                        &copy; {{.Year}} NATS Framework. All rights reserved.
                    </div>
                    <div class="flex items-center gap-4 text-sm text-gray-500 dark:text-gray-400">
                        <span>Version {{.Version}}</span>
                        <span class="w-px h-4 bg-gray-300 dark:bg-gray-600"></span>
                        <span>Go {{.GoVersion}}</span>
                    </div>
                </div>
            </div>
        </footer>
    </div>
    
    <!-- ============================================ -->
    <!-- SCRIPTS -->
    <!-- ============================================ -->
    <script>
        // إضافة مستمع للأحداث من HTMX
        document.addEventListener('htmx:afterSwap', function(evt) {
            // إعادة تهيئة Alpine.js بعد تبديل المحتوى
            if (typeof Alpine !== 'undefined') {
                Alpine.initTree(evt.detail.target);
            }
        });
        
        // إضافة مستمع للأخطاء
        document.addEventListener('htmx:responseError', function(evt) {
            const message = evt.detail.xhr?.responseText || 'An error occurred';
            try {
                const data = JSON.parse(message);
                if (data.error) {
                    // عرض الإشعار عبر Alpine
                    if (window.app) {
                        window.app.addNotification(data.error, 'danger');
                    }
                }
            } catch(e) {
                // عرض الإشعار عبر Alpine
                if (window.app) {
                    window.app.addNotification('An error occurred', 'danger');
                }
            }
        });
        
        // إضافة مستمع للطلبات الناجحة
        document.addEventListener('htmx:afterRequest', function(evt) {
            if (evt.detail.successful) {
                const trigger = evt.detail.xhr?.getResponseHeader('HX-Trigger');
                if (trigger) {
                    try {
                        const data = JSON.parse(trigger);
                        if (data.success && data.message) {
                            if (window.app) {
                                window.app.addNotification(data.message, 'success');
                            }
                        }
                    } catch(e) {}
                }
            }
        });
        
        // إضافة مستمع لإغلاق المودال
        document.addEventListener('htmx:afterSwap', function(evt) {
            if (evt.detail.target.id === 'modal-content') {
                // فتح المودال تلقائياً
                if (window.app && !window.app.showModal) {
                    window.app.openModal();
                }
            }
        });
        
        console.log('🚀 NATS Framework loaded successfully!');
        console.log(`📦 Version: {{.Version}}`);
        console.log(`🌐 Environment: {{.Env}}`);
    </script>
</body>
</html>