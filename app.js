const API_URL = '/api';

const playerBoardEl = document.getElementById('player-board');
const enemyBoardEl = document.getElementById('enemy-board');
const messageAreaEl = document.getElementById('message-area');
const newGameButton = document.getElementById('new-game-button');
const abilitiesListEl = document.getElementById('abilities-list');
const saveGameButton = document.getElementById('save-game-button');
const loadGameButton = document.getElementById('load-game-button');

async function updateGameView() {
    try {
        const response = await fetch(`${API_URL}/game`);
        if (!response.ok) {
            throw new Error('Не удалось загрузить состояние игры (сервер недоступен?)');
        }

        const data = await response.json();
        const gameState = data.game;

        renderBoard(playerBoardEl, gameState.Player1.MyBoard.Grid, false);
        renderBoard(enemyBoardEl, gameState.Player2.MyBoard.Grid, true);

        updateMessage(gameState);
        updateAbilities(gameState.Player1.Abilities);

        if (data.save_exists) {
            loadGameButton.style.display = 'inline-block';
        } else {
            loadGameButton.style.display = 'none';
        }

    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
    }
}

function renderBoard(tableElement, grid, isEnemy) {
    tableElement.innerHTML = '';
    for (let i = 0; i < 10; i++) {
        const row = document.createElement('tr');
        for (let j = 0; j < 10; j++) {
            const cell = document.createElement('td');
            const cellState = grid[i][j];
            switch (cellState) {
                case 0: cell.className = 'cell-empty'; break;
                case 1: cell.className = isEnemy ? 'cell-empty' : 'cell-ship'; break;
                case 2: cell.className = 'cell-miss'; cell.textContent = '•'; break;
                case 3: cell.className = 'cell-hit'; cell.textContent = '✕'; break;
            }
            if (isEnemy && (cellState === 0 || cellState === 1)) {
                cell.addEventListener('click', () => handleAttackClick(i, j));
            }
            row.appendChild(cell);
        }
        tableElement.appendChild(row);
    }
}

function updateMessage(gameState) {
    if (gameState.CurrentPlayer.Name === 'Player') {
        messageAreaEl.textContent = "Ваш ход.";
    } else {
        messageAreaEl.textContent = "Ход компьютера...";
    }
}

function updateAbilities(abilities) {
    abilitiesListEl.innerHTML = '';
    if (!abilities || abilities.length === 0) {
        const li = document.createElement('li');
        li.textContent = 'У вас нет способностей.';
        abilitiesListEl.appendChild(li);
        return;
    }
    abilities.forEach(ability => {
        const li = document.createElement('li');
        li.textContent = ability.Name;
        abilitiesListEl.appendChild(li);
    });
}

async function handleAttackClick(x, y) {
    try {
        messageAreaEl.textContent = `Атакуем клетку (${x}, ${y})...`;
        const response = await fetch(`${API_URL}/attack?x=${x}&y=${y}`, { method: 'POST' });

        const result = await response.json();
        if (!response.ok) throw new Error(result.Message || 'Неизвестная ошибка атаки');

        let fullMessage = result.message;
        if (result.ability_result) {
            fullMessage += ` | ${result.ability_result}`;
        }
        messageAreaEl.textContent = fullMessage;

        await updateGameView();

        if (result.game_over) {
            messageAreaEl.textContent = `Игра окончена! Победитель: ${result.winner}`;
            document.querySelectorAll('.board').forEach(b => b.style.opacity = '0.5');
        }

    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
        await updateGameView();
    }
}

async function handleNewGameClick() {
    try {
        messageAreaEl.textContent = 'Создание новой игры...';
        const response = await fetch(`${API_URL}/newgame`, { method: 'POST' });
        if (!response.ok) throw new Error('Не удалось начать новую игру');

        await updateGameView();
        messageAreaEl.textContent = "Новая игра началась. Ваш ход.";
    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
    }
}

async function handleSaveGameClick() {
    try {
        messageAreaEl.textContent = 'Сохранение игры...';
        const response = await fetch(`${API_URL}/save`, { method: 'POST' });

        const result = await response.json();
        if (!response.ok) throw new Error(result.message || 'Ошибка сохранения');

        messageAreaEl.textContent = result.message;
        await updateGameView();
    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
    }
}

async function handleLoadGameClick() {
    try {
        messageAreaEl.textContent = 'Загрузка игры...';
        const response = await fetch(`${API_URL}/load`, { method: 'POST' });

        const result = await response.json();
        if (!response.ok) throw new Error(result.message || 'Ошибка загрузки');

        await updateGameView();
        messageAreaEl.textContent = result.message;
    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
    }
}

newGameButton.addEventListener('click', handleNewGameClick);
saveGameButton.addEventListener('click', handleSaveGameClick);
loadGameButton.addEventListener('click', handleLoadGameClick);

document.addEventListener('DOMContentLoaded', updateGameView);