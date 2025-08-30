const API_URL = '/api';

const playerBoardEl = document.getElementById('player-board');
const enemyBoardEl = document.getElementById('enemy-board');
const messageAreaEl = document.getElementById('message-area');
const newGameButton = document.getElementById('new-game-button');
const abilitiesListEl = document.getElementById('abilities-list');
const saveGameButton = document.getElementById('save-game-button');
const loadGameButton = document.getElementById('load-game-button');

let isAnimating = false;
let selectedAbility = null;

async function updateGameView() {
    try {
        const response = await fetch(`${API_URL}/game`);
        if (!response.ok) throw new Error("Не удалось загрузить игру");
        const data = await response.json();
        const gameState = data.game;

        renderBoard(playerBoardEl, gameState.Player1.MyBoard.Grid, false);
        renderBoard(enemyBoardEl, gameState.Player2.MyBoard.Grid, true);
        renderAbilities(gameState.Player1.Abilities);
        updateMessage(gameState);

        loadGameButton.style.display = data.save_exists ? 'inline-block' : 'none';
        isAnimating = false;

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
                cell.addEventListener('click', () => onEnemyCellClick(i, j));
            }
            row.appendChild(cell);
        }
        tableElement.appendChild(row);
    }
}

function renderAbilities(abilities) {
    abilitiesListEl.innerHTML = '';
    if (!abilities || abilities.length === 0) {
        abilitiesListEl.innerHTML = '<p>Нет способностей</p>';
        return;
    }
    abilities.forEach(ability => {
        const button = document.createElement('button');
        button.className = 'ability-button';
        button.textContent = ability.Name;
        button.dataset.abilityName = ability.Name;
        button.dataset.requiresTarget = ability.RequiresTarget;
        button.addEventListener('click', onAbilityClick);
        abilitiesListEl.appendChild(button);
    });
}

function updateMessage(gameState) {
    if (selectedAbility) return;
    if (gameState.CurrentPlayer.Name === 'Player') {
        messageAreaEl.textContent = "Ваш ход.";
    } else {
        messageAreaEl.textContent = "Ход компьютера...";
    }
}

function onAbilityClick(event) {
    if (isAnimating) return;
    const button = event.target;
    const abilityName = button.dataset.abilityName;
    const requiresTarget = button.dataset.requiresTarget === 'true';

    if (button.classList.contains('selected')) {
        selectedAbility = null;
        button.classList.remove('selected');
        enemyBoardEl.classList.remove('targeting-mode');
        messageAreaEl.textContent = "Выбор цели отменен. Ваш ход.";
        return;
    }

    if (requiresTarget) {
        selectedAbility = { name: abilityName, button: button };
        document.querySelectorAll('.ability-button').forEach(b => b.classList.remove('selected'));
        button.classList.add('selected');
        enemyBoardEl.classList.add('targeting-mode');
        messageAreaEl.textContent = `Выберите цель для способности "${abilityName}"`;
    } else {
        useAbility(abilityName);
    }
}

function onEnemyCellClick(x, y) {
    if (isAnimating) return;
    if (selectedAbility) {
        useAbility(selectedAbility.name, x, y);
    } else {
        handleAttack(x, y);
    }
}

async function handleAttack(x, y) {
    isAnimating = true;
    messageAreaEl.textContent = `Атакуем клетку (${x}, ${y})...`;
    try {
        const response = await fetch(`${API_URL}/attack?x=${x}&y=${y}`, { method: 'POST' });
        const result = await response.json();
        if (!response.ok) throw new Error(result.Message || 'Ошибка атаки');
        if (result.human_move_result) {
            await animateMove(enemyBoardEl, result.human_move_result);
        }
        if (result.game_over) {
            handleGameOver(result.winner);
            return;
        }
        messageAreaEl.textContent = result.message;
        if (result.computer_moves && result.computer_moves.length > 0) {
            messageAreaEl.textContent = "Ход компьютера...";
            for (const move of result.computer_moves) {
                await sleep(800);
                await animateMove(playerBoardEl, move);
            }
        }
        await updateGameView();
    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
        isAnimating = false;
    }
}

async function useAbility(abilityName, x, y) {
    isAnimating = true;
    let url = `${API_URL}/ability?ability_name=${abilityName}`;
    if (x !== undefined && y !== undefined) {
        url += `&x=${x}&y=${y}`;
    }
    if (selectedAbility) {
        selectedAbility.button.classList.remove('selected');
        enemyBoardEl.classList.remove('targeting-mode');
        selectedAbility = null;
    }
    try {
        const response = await fetch(url, { method: 'POST' });
        const result = await response.json();
        if (!response.ok) throw new Error(result.Message || 'Ошибка способности');

        messageAreaEl.textContent = result.message;
        if (result.attack_result) {
            await animateMove(enemyBoardEl, result.attack_result);
        }
        if (result.affected_points) {
            await animateScan(enemyBoardEl, result.affected_points);
        }

        if (result.game_over) {
            handleGameOver(result.winner);
            return;
        }
        await updateGameView();
    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
        isAnimating = false;
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function animateScan(boardElement, points) {
    if (!points || points.length === 0) return;

    for (const p of points) {
        const cell = boardElement.rows[p.X].cells[p.Y];
        cell.classList.remove('cell-scanned');
        void cell.offsetWidth;
        cell.classList.add('cell-scanned');
    }

    await sleep(4000);
}

async function animateMove(boardElement, moveData) {
    if (!moveData) return;

    const move = {
        x: moveData.x ?? moveData.Target?.X,
        y: moveData.y ?? moveData.Target?.Y,
        result: moveData.result ?? moveData.Result,
        marked_points: moveData.marked_points ?? moveData.MarkedPoints
    };

    if (move.x === undefined || move.y === undefined) {
        console.error("Не удалось определить координаты хода из данных:", moveData);
        return;
    }

    const cell = boardElement.rows[move.x].cells[move.y];
    cell.classList.add('cell-attacked');
    await sleep(200);

    if (move.result === 0) {
        cell.className = 'cell-miss'; cell.textContent = '•';
    } else {
        cell.className = 'cell-hit'; cell.textContent = '✕';
    }

    await sleep(200);

    if (move.result === 2 && move.marked_points) {
        messageAreaEl.textContent = "Потопил!";
        for (const p of move.marked_points) {
            const markedCell = boardElement.rows[p.X].cells[p.Y];
            if (markedCell.className.includes('empty') || markedCell.className.includes('ship')) {
                markedCell.className = 'cell-miss';
                markedCell.textContent = '•';
                await sleep(50);
            }
        }
    }
}

function handleGameOver(winner) {
    messageAreaEl.textContent = `Игра окончена! Победитель: ${winner}`;
    isAnimating = true;
}

newGameButton.addEventListener('click', async () => {
    if (isAnimating) return;
    await fetch(`${API_URL}/newgame`, { method: 'POST' });
    updateGameView();
});
saveGameButton.addEventListener('click', () => {
    if (isAnimating) return;
    fetch(`${API_URL}/save`, { method: 'POST' });
});
loadGameButton.addEventListener('click', async () => {
    if (isAnimating) return;
    await fetch(`${API_URL}/load`, { method: 'POST' });
    updateGameView();
});

document.addEventListener('DOMContentLoaded', updateGameView);