const API_URL = '/api';

const playerBoardEl = document.getElementById('player-board');
const enemyBoardEl = document.getElementById('enemy-board');
const messageAreaEl = document.getElementById('message-area');
const newGameButton = document.getElementById('new-game-button');
const abilitiesListEl = document.getElementById('abilities-list');
const saveGameButton = document.getElementById('save-game-button');
const loadGameButton = document.getElementById('load-game-button');

let isAnimating = false;

async function updateGameView() {
    try {
        const response = await fetch(`${API_URL}/game`);
        const data = await response.json();
        if (!response.ok) throw new Error("Не удалось загрузить игру");
        const gameState = data.game;

        renderBoard(playerBoardEl, gameState.Player1.MyBoard.Grid, false);
        renderBoard(enemyBoardEl, gameState.Player2.MyBoard.Grid, true);
        updateMessage(gameState);
        updateAbilities(gameState.Player1.Abilities);

        loadGameButton.style.display = data.save_exists ? 'inline-block' : 'none';
        document.querySelectorAll('.board').forEach(b => b.style.opacity = '1');
        isAnimating = false;

        if (gameState.CurrentPlayer.Name === 'Computer') {
            await handleComputerMoves();
        }

    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
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
    if (isAnimating) return;
    isAnimating = true;
    messageAreaEl.textContent = `Атакуем клетку (${x}, ${y})...`;

    try {
        const response = await fetch(`${API_URL}/attack?x=${x}&y=${y}`, { method: 'POST' });
        const result = await response.json();
        if (!response.ok) throw new Error(result.Message || 'Неизвестная ошибка атаки');

        if (result.human_move_result) {
            await animateMove(enemyBoardEl, result.human_move_result, false);
        }

        if (result.computer_moves && result.computer_moves.length > 0) {
            for (const move of result.computer_moves) {
                await sleep(800);
                await animateMove(playerBoardEl, move, true);
            }
        }

        const finalStateResponse = await fetch(`${API_URL}/game`);
        const finalData = await finalStateResponse.json();
        updateMessage(finalData.game);
        updateAbilities(finalData.game.Player1.Abilities);

        messageAreaEl.textContent = result.message;

        if (result.game_over) {
            handleGameOver(result.winner);
        } else {
            isAnimating = false;
        }

    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
        isAnimating = false;
    }
}

async function animateMove(boardElement, move, isPlayerBoard) {
    if (!move) return;
    const cell = boardElement.rows[move.x].cells[move.y];

    cell.replaceWith(cell.cloneNode(true));
    const newCell = boardElement.rows[move.x].cells[move.y];

    const highlightClass = isPlayerBoard ? 'cell-attacked-by-bot' : 'cell-attacked-by-player';
    newCell.classList.add(highlightClass);
    messageAreaEl.textContent = `${isPlayerBoard ? 'Компьютер' : 'Вы'} атакует (${move.x}, ${move.y})...`;
    await sleep(400);

    if (move.result === 0) {
        newCell.className = 'cell-miss';
        newCell.textContent = '•';
    } else {
        newCell.className = 'cell-hit';
        newCell.textContent = '✕';
    }

    await sleep(600);

    if (move.result === 2 && move.marked_points) {
        messageAreaEl.textContent = "Потопил!";
        for (const p of move.marked_points) {
            const markedCell = boardElement.rows[p.X].cells[p.Y];
            if (markedCell.className === 'cell-empty' || markedCell.className === 'cell-ship') {
                markedCell.className = 'cell-miss';
                markedCell.textContent = '•';
                await sleep(50);
            }
        }
    }
}

function handleGameOver(winner) {
    messageAreaEl.textContent = `Игра окончена! Победитель: ${winner}`;
    document.querySelectorAll('.board').forEach(b => b.style.opacity = '0.5');
    isAnimating = true;
}

async function handleNewGameClick() {
    if (isAnimating) return;
    isAnimating = true;
    const response = await fetch(`${API_URL}/newgame`, { method: 'POST' });
    if (response.ok) await updateGameView();
}

async function handleSaveGameClick() {
    if (isAnimating) return;
    const response = await fetch(`${API_URL}/save`, { method: 'POST' });
    if (response.ok) await updateGameView();
}

async function handleLoadGameClick() {
    if (isAnimating) return;
    isAnimating = true;
    const response = await fetch(`${API_URL}/load`, { method: 'POST' });
    if (response.ok) await updateGameView();
}

newGameButton.addEventListener('click', handleNewGameClick);
saveGameButton.addEventListener('click', handleSaveGameClick);
loadGameButton.addEventListener('click', handleLoadGameClick);

document.addEventListener('DOMContentLoaded', updateGameView);