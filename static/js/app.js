/* ============================================
   NATS Framework - Main JavaScript (Local)
   ============================================ */

// No CDN - All local

// HTMX is loaded from local file
// Alpine.js is loaded from local file

console.log('🚀 NATS Framework loaded successfully!');

// HTMX Configuration
document.addEventListener('DOMContentLoaded', function() {
    if (typeof htmx !== 'undefined') {
        htmx.config.historyCacheSize = 20;
        htmx.config.defaultSwapStyle = 'innerHTML';
        htmx.config.defaultSwapDelay = 0;
        htmx.config.defaultTimeout = 30000;
        console.log('✅ HTMX configured (local)');
    }
    
    if (typeof Alpine !== 'undefined') {
        console.log('✅ Alpine.js configured (local)');
    }
});

// Toast notifications
function showToast(message, type) {
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

// Modal helpers
function openModal(id) {
    var modal = document.getElementById(id);
    if (modal) {
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }
}

function closeModal(id) {
    var modal = document.getElementById(id);
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = 'auto';
    }
}

// Alpine.js App for Developer Dashboard
function developerApp() {
    return {
        mode: 'developer',
        currentView: 'dashboard',
        pageTitle: '📊 Dashboard',
        loading: false,
        creating: false,
        linking: false,
        creatingUser: false,
        addingField: false,
        pendingMigration: false,
        
        stats: { apps: 0, models: 0, users: 0, relations: 0 },
        apps: [],
        selectedApp: null,
        selectedFields: [],
        relations: [],
        
        showCreateApp: false,
        newApp: {
            name: '',
            description: '',
            parent: 'core',
            models: [{name: '', fields: ''}],
            crud: 'yes'
        },
        
        showLinkApp: false,
        linkApp: {
            parent: '',
            child: '',
            type: 'one-to-many',
            foreignKey: ''
        },
        
        showCreateUser: false,
        newUser: {
            username: '',
            email: '',
            fullName: '',
            password: '',
            role: 'user'
        },
        
        showAddField: false,
        newField: {
            name: '',
            type: 'string',
            required: 'false',
            unique: 'false',
            default: '',
            relatedModel: ''
        },
        
        init() {
            console.log('🛠️ Developer Dashboard initialized');
            this.loadStats();
            this.loadApps();
            this.loadRelations();
            window.dashboardAppInstance = this;
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
                .catch(err => {
                    console.error('Apps error:', err);
                    this.showToast('Failed to load apps', 'error');
                });
        },
        
        loadRelations() {
            fetch('/api/admin/relations')
                .then(res => res.json())
                .then(data => {
                    if (data.success) this.relations = data.data;
                })
                .catch(err => {
                    console.error('Relations error:', err);
                    this.showToast('Failed to load relations', 'error');
                });
        },
        
        viewApp(appName) {
            this.selectedApp = appName;
            this.pendingMigration = false;
            fetch('/api/admin/apps/' + appName + '/models')
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        this.selectedFields = data.data;
                    } else {
                        this.selectedFields = [];
                        this.showToast('No models found for ' + appName, 'info');
                    }
                })
                .catch(err => {
                    console.error('Models error:', err);
                    this.showToast('Failed to load models', 'error');
                });
        },
        
        openCreateApp() {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            this.newApp = { name: '', description: '', parent: 'core', models: [{name: '', fields: ''}], crud: 'yes' };
            this.showCreateApp = true;
        },
        
        executeCreateApp() {
            if (!this.newApp.name) {
                this.showToast('Please enter application name', 'error');
                return;
            }
            this.creating = true;
            fetch('/api/admin/apps', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.newApp)
            })
            .then(res => res.json())
            .then(data => {
                this.creating = false;
                if (data.success) {
                    this.showCreateApp = false;
                    this.loadApps();
                    this.showToast('✅ Application "' + this.newApp.name + '" created successfully!', 'success');
                } else {
                    this.showToast('❌ Error: ' + (data.error || 'Unknown error'), 'error');
                }
            })
            .catch(err => {
                this.creating = false;
                this.showToast('❌ Network error: ' + err.message, 'error');
            });
        },
        
        openLinkApp(appName) {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            this.linkApp = { parent: appName || '', child: '', type: 'one-to-many', foreignKey: '' };
            this.showLinkApp = true;
        },
        
        executeLinkApp() {
            if (!this.linkApp.parent || !this.linkApp.child) {
                this.showToast('Please select both parent and child apps', 'error');
                return;
            }
            if (this.linkApp.parent === this.linkApp.child) {
                this.showToast('Parent and child apps cannot be the same', 'error');
                return;
            }
            this.linking = true;
            fetch('/api/admin/relations', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.linkApp)
            })
            .then(res => res.json())
            .then(data => {
                this.linking = false;
                if (data.success) {
                    this.showLinkApp = false;
                    this.loadRelations();
                    this.showToast('✅ Apps linked successfully!', 'success');
                } else {
                    this.showToast('❌ Error: ' + (data.error || 'Unknown error'), 'error');
                }
            })
            .catch(err => {
                this.linking = false;
                this.showToast('❌ Network error: ' + err.message, 'error');
            });
        },
        
        openAddField() {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            if (!this.selectedApp) {
                this.showToast('Please select an app first', 'error');
                return;
            }
            this.newField = { name: '', type: 'string', required: 'false', unique: 'false', default: '', relatedModel: '' };
            this.showAddField = true;
        },
        
        executeAddField() {
            if (!this.newField.name) {
                this.showToast('Please enter a field name', 'error');
                return;
            }
            this.addingField = true;
            fetch('/api/admin/apps/' + this.selectedApp + '/models/Default/fields', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    name: this.newField.name,
                    type: this.newField.type,
                    required: this.newField.required === 'true',
                    unique: this.newField.unique === 'true',
                    default: this.newField.default,
                    relatedModel: this.newField.relatedModel
                })
            })
            .then(res => res.json())
            .then(data => {
                this.addingField = false;
                if (data.success) {
                    this.showAddField = false;
                    this.showToast('✅ Field "' + this.newField.name + '" added successfully!', 'success');
                    this.viewApp(this.selectedApp);
                    this.pendingMigration = true;
                    this.showMigrationNotification();
                } else {
                    this.showToast('❌ Error: ' + data.error, 'error');
                }
            })
            .catch(err => {
                this.addingField = false;
                this.showToast('❌ Network error: ' + err.message, 'error');
            });
        },
        
        deleteApp(appName) {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            if (appName === 'core' || appName === 'users') {
                this.showToast('❌ Cannot delete system apps (core, users)', 'error');
                return;
            }
            if (!confirm('Are you sure you want to delete app: "' + appName + '"?')) return;
            fetch('/api/admin/apps/' + appName, { method: 'DELETE' })
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        this.loadApps();
                        this.loadStats();
                        this.selectedApp = null;
                        this.selectedFields = [];
                        this.showToast('✅ App "' + appName + '" deleted successfully!', 'success');
                    } else {
                        this.showToast('❌ Error: ' + (data.error || 'Unknown error'), 'error');
                    }
                })
                .catch(err => {
                    this.showToast('❌ Network error: ' + err.message, 'error');
                });
        },
        
        deleteRelation(rel) {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            if (!confirm('Delete relation between "' + rel.parent + '" and "' + rel.child + '"?')) return;
            fetch('/api/admin/relations', {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(rel)
            })
            .then(res => res.json())
            .then(data => {
                if (data.success) {
                    this.loadRelations();
                    this.loadStats();
                    this.showToast('✅ Relation deleted successfully!', 'success');
                } else {
                    this.showToast('❌ Error: ' + (data.error || 'Unknown error'), 'error');
                }
            })
            .catch(err => {
                this.showToast('❌ Network error: ' + err.message, 'error');
            });
        },
        
        runMigrations() {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            this.showToast('⏳ Running migrations...', 'info');
            this.pendingMigration = false;
            fetch('/api/admin/migrations/run', { method: 'POST' })
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        this.showToast('✅ Migrations completed successfully!', 'success');
                        this.viewApp(this.selectedApp);
                    } else {
                        this.showToast('❌ Error: ' + (data.error || 'Unknown error'), 'error');
                    }
                })
                .catch(err => {
                    this.showToast('❌ Network error: ' + err.message, 'error');
                });
        },
        
        resetMigrations() {
            if (this.mode !== 'developer') {
                this.showToast('Please switch to Developer Mode first', 'error');
                return;
            }
            if (!confirm('⚠️ This will reset all migrations. Are you sure?')) return;
            this.showToast('⏳ Resetting migrations...', 'info');
            fetch('/api/admin/migrations/reset', { method: 'POST' })
                .then(res => res.json())
                .then(data => {
                    if (data.success) {
                        this.showToast('✅ Migrations reset successfully!', 'success');
                    } else {
                        this.showToast('❌ Error: ' + (data.error || 'Unknown error'), 'error');
                    }
                })
                .catch(err => {
                    this.showToast('❌ Network error: ' + err.message, 'error');
                });
        },
        
        showMigrationNotification() {
            var container = document.getElementById('toast-container');
            if (!container) return;
            var toast = document.createElement('div');
            toast.className = 'toast toast-warning';
            toast.style.background = '#fefcbf';
            toast.style.color = '#744210';
            toast.style.border = '1px solid #ecc94b';
            toast.innerHTML = '<div style="display: flex; align-items: center; gap: 12px; justify-content: space-between; width: 100%;">' +
                '<div style="display: flex; align-items: center; gap: 8px;">' +
                    '<span style="font-size: 20px;">⚠️</span>' +
                    '<span><strong>Database changes detected!</strong> Run migration to apply changes.</span>' +
                '</div>' +
                '<button onclick="this.parentElement.parentElement.remove(); window.dashboardAppInstance.runMigrations();" ' +
                        'style="background: #d69e2e; border: none; padding: 6px 16px; border-radius: 6px; cursor: pointer; color: white; font-weight: 500;">' +
                    '▶️ Run Migration' +
                '</button>' +
            '</div>';
            container.appendChild(toast);
            setTimeout(function() {
                toast.style.opacity = '0';
                toast.style.transition = 'opacity 0.3s';
                setTimeout(function() { toast.remove(); }, 300);
            }, 15000);
        },
        
        showToast: function(message, type) {
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