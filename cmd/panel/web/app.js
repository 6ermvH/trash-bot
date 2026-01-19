const API_BASE = '/api';

const state = {
    token: localStorage.getItem('token') || null
};

const elements = {
    loginSection: document.getElementById('login-section'),
    dashboardSection: document.getElementById('dashboard-section'),
    loginForm: document.getElementById('login-form'),
    loginError: document.getElementById('login-error'),
    logoutBtn: document.getElementById('logout-btn'),
    totalChats: document.getElementById('total-chats'),
    totalUsers: document.getElementById('total-users'),
    avgUsers: document.getElementById('avg-users'),
    chatsBody: document.getElementById('chats-body'),
    chatsTable: document.getElementById('chats-table'),
    noChats: document.getElementById('no-chats')
};

async function apiRequest(endpoint, options = {}) {
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (state.token) {
        headers['Authorization'] = `Bearer ${state.token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers
    });

    if (response.status === 401) {
        logout();
        throw new Error('Unauthorized');
    }

    return response;
}

async function login(loginValue, password) {
    const response = await fetch(`${API_BASE}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login: loginValue, password })
    });

    if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Login failed');
    }

    const data = await response.json();
    state.token = data.token;
    localStorage.setItem('token', data.token);
}

function logout() {
    state.token = null;
    localStorage.removeItem('token');
    showLogin();
}

function showLogin() {
    elements.loginSection.classList.remove('hidden');
    elements.dashboardSection.classList.add('hidden');
}

function showDashboard() {
    elements.loginSection.classList.add('hidden');
    elements.dashboardSection.classList.remove('hidden');
    loadDashboard();
}

async function loadDashboard() {
    await Promise.all([loadStats(), loadChats()]);
}

async function loadStats() {
    try {
        const response = await apiRequest('/stats');
        const stats = await response.json();

        elements.totalChats.textContent = stats.totalChats;
        elements.totalUsers.textContent = stats.totalUsers;
        elements.avgUsers.textContent = stats.avgUsersPerChat.toFixed(1);
    } catch (error) {
        console.error('Failed to load stats:', error);
    }
}

async function loadChats() {
    try {
        const response = await apiRequest('/chats');
        const chats = await response.json();

        if (chats.length === 0) {
            elements.chatsTable.classList.add('hidden');
            elements.noChats.classList.remove('hidden');
            return;
        }

        elements.chatsTable.classList.remove('hidden');
        elements.noChats.classList.add('hidden');

        elements.chatsBody.innerHTML = chats.map(chat => `
            <tr>
                <td>${chat.id}</td>
                <td>${chat.activeUsers[chat.currentUser] || '-'}</td>
                <td>${chat.activeUsers.join(', ')}</td>
            </tr>
        `).join('');
    } catch (error) {
        console.error('Failed to load chats:', error);
    }
}

elements.loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    elements.loginError.textContent = '';

    const loginValue = document.getElementById('login').value;
    const password = document.getElementById('password').value;

    try {
        await login(loginValue, password);
        showDashboard();
    } catch (error) {
        elements.loginError.textContent = error.message;
    }
});

elements.logoutBtn.addEventListener('click', logout);

// Initial state
if (state.token) {
    showDashboard();
} else {
    showLogin();
}
