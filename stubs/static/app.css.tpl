/* ============================================
   NATS Framework - Main Stylesheet
   ============================================ */

/* ============================================
   CSS Variables
   ============================================ */
:root {
    --primary-color: #4F46E5;
    --primary-hover: #4338CA;
    --secondary-color: #6B7280;
    --success-color: #10B981;
    --danger-color: #EF4444;
    --warning-color: #F59E0B;
    --info-color: #3B82F6;
    
    --text-primary: #111827;
    --text-secondary: #6B7280;
    --text-light: #9CA3AF;
    
    --bg-primary: #FFFFFF;
    --bg-secondary: #F3F4F6;
    --bg-tertiary: #E5E7EB;
    
    --border-color: #D1D5DB;
    --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    --shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    
    --radius: 0.5rem;
    --transition: 0.2s ease-in-out;
}

/* ============================================
   Reset & Base
   ============================================ */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    color: var(--text-primary);
    background-color: var(--bg-secondary);
    line-height: 1.6;
}

/* ============================================
   Layout
   ============================================ */
.container {
    max-width: 1280px;
    margin: 0 auto;
    padding: 0 1rem;
}

/* ============================================
   Navigation
   ============================================ */
.navbar {
    background-color: var(--bg-primary);
    border-bottom: 1px solid var(--border-color);
    padding: 0.75rem 1rem;
    box-shadow: var(--shadow-sm);
    position: sticky;
    top: 0;
    z-index: 50;
}

.navbar-brand {
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--primary-color);
    text-decoration: none;
}

/* ============================================
   Cards
   ============================================ */
.card {
    background-color: var(--bg-primary);
    border-radius: var(--radius);
    box-shadow: var(--shadow);
    padding: 1.5rem;
    margin-bottom: 1rem;
    transition: box-shadow var(--transition);
}

.card:hover {
    box-shadow: var(--shadow-md);
}

.card-header {
    font-size: 1.125rem;
    font-weight: 600;
    margin-bottom: 1rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--border-color);
}

/* ============================================
   Buttons
   ============================================ */
.btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    font-weight: 500;
    border-radius: var(--radius);
    border: 1px solid transparent;
    cursor: pointer;
    transition: all var(--transition);
    text-decoration: none;
    gap: 0.5rem;
}

.btn-primary {
    background-color: var(--primary-color);
    color: white;
}

.btn-primary:hover {
    background-color: var(--primary-hover);
}

.btn-secondary {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
}

.btn-secondary:hover {
    background-color: var(--border-color);
}

.btn-success {
    background-color: var(--success-color);
    color: white;
}

.btn-success:hover {
    background-color: #059669;
}

.btn-danger {
    background-color: var(--danger-color);
    color: white;
}

.btn-danger:hover {
    background-color: #DC2626;
}

.btn-sm {
    padding: 0.25rem 0.5rem;
    font-size: 0.75rem;
}

.btn-lg {
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
}

/* ============================================
   Tables
   ============================================ */
.table-container {
    overflow-x: auto;
    border-radius: var(--radius);
    border: 1px solid var(--border-color);
}

.table {
    width: 100%;
    border-collapse: collapse;
    background-color: var(--bg-primary);
}

.table thead {
    background-color: var(--bg-secondary);
}

.table th {
    padding: 0.75rem 1rem;
    text-align: left;
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-secondary);
    border-bottom: 1px solid var(--border-color);
}

.table td {
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--border-color);
}

.table tbody tr:hover {
    background-color: var(--bg-secondary);
}

/* ============================================
   Forms
   ============================================ */
.form-group {
    margin-bottom: 1rem;
}

.form-label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 0.25rem;
}

.form-control {
    width: 100%;
    padding: 0.5rem 0.75rem;
    font-size: 0.875rem;
    border: 1px solid var(--border-color);
    border-radius: var(--radius);
    transition: border-color var(--transition), box-shadow var(--transition);
    background-color: var(--bg-primary);
}

.form-control:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.form-control.is-invalid {
    border-color: var(--danger-color);
}

.form-control.is-valid {
    border-color: var(--success-color);
}

/* ============================================
   Alerts
   ============================================ */
.alert {
    padding: 0.75rem 1rem;
    border-radius: var(--radius);
    margin-bottom: 1rem;
}

.alert-success {
    background-color: #D1FAE5;
    color: #065F46;
    border: 1px solid #A7F3D0;
}

.alert-danger {
    background-color: #FEE2E2;
    color: #991B1B;
    border: 1px solid #FCA5A5;
}

.alert-warning {
    background-color: #FEF3C7;
    color: #92400E;
    border: 1px solid #FCD34D;
}

.alert-info {
    background-color: #DBEAFE;
    color: #1E40AF;
    border: 1px solid #93C5FD;
}

/* ============================================
   Badges
   ============================================ */
.badge {
    display: inline-block;
    padding: 0.25rem 0.75rem;
    font-size: 0.75rem;
    font-weight: 600;
    border-radius: 9999px;
}

.badge-primary {
    background-color: #DBEAFE;
    color: #1E40AF;
}

.badge-success {
    background-color: #D1FAE5;
    color: #065F46;
}

.badge-danger {
    background-color: #FEE2E2;
    color: #991B1B;
}

.badge-warning {
    background-color: #FEF3C7;
    color: #92400E;
}

/* ============================================
   Utilities
   ============================================ */
.text-center { text-align: center; }
.text-right { text-align: right; }
.text-left { text-align: left; }

.mt-1 { margin-top: 0.25rem; }
.mt-2 { margin-top: 0.5rem; }
.mt-3 { margin-top: 0.75rem; }
.mt-4 { margin-top: 1rem; }
.mt-5 { margin-top: 1.5rem; }

.mb-1 { margin-bottom: 0.25rem; }
.mb-2 { margin-bottom: 0.5rem; }
.mb-3 { margin-bottom: 0.75rem; }
.mb-4 { margin-bottom: 1rem; }
.mb-5 { margin-bottom: 1.5rem; }

.p-1 { padding: 0.25rem; }
.p-2 { padding: 0.5rem; }
.p-3 { padding: 0.75rem; }
.p-4 { padding: 1rem; }
.p-5 { padding: 1.5rem; }

.d-flex { display: flex; }
.d-grid { display: grid; }
.d-none { display: none; }

.flex-column { flex-direction: column; }
.align-items-center { align-items: center; }
.justify-content-between { justify-content: space-between; }
.gap-1 { gap: 0.25rem; }
.gap-2 { gap: 0.5rem; }
.gap-3 { gap: 0.75rem; }

.w-100 { width: 100%; }
.h-100 { height: 100%; }

/* ============================================
   Responsive
   ============================================ */
@media (max-width: 768px) {
    .container {
        padding: 0 0.5rem;
    }
    
    .table {
        font-size: 0.75rem;
    }
    
    .table th,
    .table td {
        padding: 0.5rem;
    }
    
    .card {
        padding: 1rem;
    }
}