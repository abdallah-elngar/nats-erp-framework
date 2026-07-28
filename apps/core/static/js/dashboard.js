// ============================================
// Developer Dashboard
// ============================================

function developerApp() {
    return {
        mode: 'developer',
        currentView: 'dashboard',
        pageTitle: '📊 Dashboard',

        stats: { apps: 0, models: 0, users: 0, relations: 0 },
        apps: [],
        relations: [],

        showCreateApp: false,
        newApp: { name: '', description: '', parent: 'core', models: [{name: '', fields: ''}], crud: 'yes' },

        showLinkApp: false,
        linkApp: { parent: '', child: '', type: 'one-to-many', foreignKey: '' },

        showAddField: false,
        newField: { name: '', type: 'string', required: false, unique: false, default: '', relatedModel: '' },

        init() {
            console.log('🚀 Developer Dashboard initialized');
            this.loadStats();
            this.loadApps();
        },

        setMode(mode) {
            this.mode = mode;
            this.showToast('Switched to ' + (mode === 'developer' ? 'Developer' : 'Production') + ' Mode', 'info');
        },

        loadStats() {
            fetch('/api/admin/stats')
                .then(res => res.json())
                .then(data => {
                    if (data.success) this.stats = data.data;
                })
                .catch(err => console.error('Stats error:', err));
        },

        loadApps() {
            fetch('/api/admin/apps')
                .then(res => res.json())
                .then(data => {
                    if (data.success) this.apps = data.data;
                })
                .catch(err => console.error('Apps error:', err));
        },

        executeCreateApp() {
            if (!this.newApp.name) {
                this.showToast('Please enter application name', 'error');
                return;
            }
            fetch('/api/admin/apps', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.newApp)
            })
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    this.showCreateApp = false;
                    this.loadApps();
                    this.showToast('✅ Application created!', 'success');
                } else {
                    this.showToast('❌ Error: ' + data.error, 'error');
                }
            })
            .catch(err => {
                this.showToast('❌ Network error: ' + err.message, 'error');
            });
        },

        showToast(message, type) {
            var container = document.getElementById('toast-container');
            if (!container) {
                container = document.createElement('div');
                container.id = 'toast-container';
                container.style.cssText = 'position: fixed; bottom: 20px; right: 20px; z-index: 2000; display: flex; flex-direction: column; gap: 8px;';
                document.body.appendChild(container);
            }
            var toast = document.createElement('div');
            toast.className = 'toast toast-' + type;
            toast.textContent = message;
            container.appendChild(toast);
            setTimeout(function() {
                toast.style.opacity = '0';
                toast.style.transition = 'opacity 0.3s';
                setTimeout(function() { toast.remove(); }, 300);
            }, 4000);
        }
    };
}