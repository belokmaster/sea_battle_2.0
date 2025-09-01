const API_URL = '/api';

const playerBoardEl = document.getElementById('player-board');
const enemyBoardEl = document.getElementById('enemy-board');
const messageAreaEl = document.getElementById('message-area');
const newGameButton = document.getElementById('new-game-button');
const abilitiesListEl = document.getElementById('abilities-list');
const saveGameButton = document.getElementById('save-game-button');
const loadGameButton = document.getElementById('load-game-button');

const newGameModal = document.getElementById('new-game-modal');
const mainGameContainer = document.getElementById('main-game-container');
const placementContainer = document.getElementById('placement-container');
const autoPlaceButton = document.getElementById('auto-place-button');
const manualPlaceButton = document.getElementById('manual-place-button');
const cancelNewGameButton = document.getElementById('cancel-new-game-button');
const placementBoardEl = document.getElementById('placement-board');
const shipListEl = document.getElementById('ship-list');
const rotateShipButton = document.getElementById('rotate-ship-button');
const resetPlacementButton = document.getElementById('reset-placement-button');
const startManualGameButton = document.getElementById('start-manual-game-button');

let isAnimating = false;
let selectedAbility = null;
let shipsToPlace = [];
let placedShips = [];
let selectedShipToPlace = null;
let isShipVertical = false;

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

function initializePlacementState() {
    shipsToPlace = [
        { id: 0, size: 4, name: 'Линкор (4)' },
        { id: 1, size: 3, name: 'Крейсер 1 (3)' },
        { id: 2, size: 3, name: 'Крейсер 2 (3)' },
        { id: 3, size: 2, name: 'Эсминец 1 (2)' },
        { id: 4, size: 2, name: 'Эсминец 2 (2)' },
        { id: 5, size: 2, name: 'Эсминец 3 (2)' },
        { id: 6, size: 1, name: 'Катер 1 (1)' },
        { id: 7, size: 1, name: 'Катер 2 (1)' },
        { id: 8, size: 1, name: 'Катер 3 (1)' },
        { id: 9, size: 1, name: 'Катер 4 (1)' }
    ];
    placedShips = [];
    selectedShipToPlace = null;
    isShipVertical = false;
    renderPlacementBoard();
    renderShipList();
    updatePlacementControls();
}

function renderPlacementBoard() {
    placementBoardEl.innerHTML = '';
    let grid = Array(10).fill(0).map(() => Array(10).fill(0));

    placedShips.forEach(ship => {
        ship.Position.forEach(p => {
            grid[p.X][p.Y] = 1;
        });
    });

    for (let i = 0; i < 10; i++) {
        const row = document.createElement('tr');
        for (let j = 0; j < 10; j++) {
            const cell = document.createElement('td');
            if (grid[i][j] === 1) {
                cell.className = 'cell-ship';
                cell.style.backgroundColor = '#777';
            } else {
                cell.className = 'cell-empty';
                cell.style.backgroundColor = '#eef7ff';
            }

            cell.dataset.x = i;
            cell.dataset.y = j;
            cell.addEventListener('mouseover', onPlacementCellMouseover);
            cell.addEventListener('mouseout', onPlacementCellMouseout);
            cell.addEventListener('click', onPlacementCellClick);
            row.appendChild(cell);
        }
        placementBoardEl.appendChild(row);
    }
}

function renderShipList() {
    shipListEl.innerHTML = '';
    shipsToPlace.forEach(ship => {
        const li = document.createElement('li');
        li.textContent = ship.name;
        li.dataset.shipId = ship.id;
        if (selectedShipToPlace && selectedShipToPlace.id === ship.id) {
            li.classList.add('selected');
        }
        li.addEventListener('click', () => {
            selectedShipToPlace = ship;
            renderShipList();
        });
        shipListEl.appendChild(li);
    });
}

function getShipPointsAndValidation(startX, startY, size, isVertical) {
    let points = [];
    let isValid = true;
    for (let i = 0; i < size; i++) {
        const x = isVertical ? startX + i : startX;
        const y = isVertical ? startY : startY + i;
        points.push({ x, y });
        if (x >= 10 || y >= 10) isValid = false;
    }
    if (!isValid) return { points, isValid };

    for (const p of points) {
        for (let dx = -1; dx <= 1; dx++) {
            for (let dy = -1; dy <= 1; dy++) {
                const checkX = p.x + dx;
                const checkY = p.y + dy;
                if (checkX >= 0 && checkX < 10 && checkY >= 0 && checkY < 10) {
                    if (placementBoardEl.rows[checkX].cells[checkY].classList.contains('cell-ship')) {
                        isValid = false;
                        break;
                    }
                }
            }
        }
        if (!isValid) break;
    }
    return { points, isValid };
}

function onPlacementCellMouseover(e) {
    if (!selectedShipToPlace) return;
    const x = parseInt(e.target.dataset.x);
    const y = parseInt(e.target.dataset.y);
    const { points, isValid } = getShipPointsAndValidation(x, y, selectedShipToPlace.size, isShipVertical);
    points.forEach(p => {
        const cell = placementBoardEl.rows[p.x]?.cells[p.y];
        if (cell) cell.classList.add(isValid ? 'cell-preview-valid' : 'cell-preview-invalid');
    });
}

function onPlacementCellMouseout(e) {
    if (!selectedShipToPlace) return;
    const x = parseInt(e.target.dataset.x);
    const y = parseInt(e.target.dataset.y);
    const { points } = getShipPointsAndValidation(x, y, selectedShipToPlace.size, isShipVertical);
    points.forEach(p => {
        const cell = placementBoardEl.rows[p.x]?.cells[p.y];
        if (cell) cell.classList.remove('cell-preview-valid', 'cell-preview-invalid');
    });
}

function onPlacementCellClick(e) {
    if (!selectedShipToPlace) return;
    const x = parseInt(e.target.dataset.x);
    const y = parseInt(e.target.dataset.y);
    const { points, isValid } = getShipPointsAndValidation(x, y, selectedShipToPlace.size, isShipVertical);

    if (isValid) {
        placedShips.push({
            Size: selectedShipToPlace.size,
            IsVertical: isShipVertical,
            Position: points.map(p => ({ X: p.x, Y: p.y })),
        });
        shipsToPlace = shipsToPlace.filter(s => s.id !== selectedShipToPlace.id);
        selectedShipToPlace = shipsToPlace[0] || null;
        renderPlacementBoard();
        renderShipList();
        updatePlacementControls();
    }
}

function updatePlacementControls() {
    rotateShipButton.textContent = `Повернуть (${isShipVertical ? 'Горизонтально' : 'Вертикально'})`;
    startManualGameButton.disabled = shipsToPlace.length > 0;
}

newGameButton.addEventListener('click', () => {
    newGameModal.style.display = 'flex';
});

cancelNewGameButton.addEventListener('click', () => {
    newGameModal.style.display = 'none';
});

autoPlaceButton.addEventListener('click', async () => {
    newGameModal.style.display = 'none';
    mainGameContainer.style.display = 'flex';
    placementContainer.style.display = 'none';
    await fetch(`${API_URL}/newgame/auto`, { method: 'POST' });
    await updateGameView();
});

manualPlaceButton.addEventListener('click', () => {
    newGameModal.style.display = 'none';
    mainGameContainer.style.display = 'none';
    placementContainer.style.display = 'flex';
    initializePlacementState();
});

rotateShipButton.addEventListener('click', () => {
    isShipVertical = !isShipVertical;
    updatePlacementControls();
});

resetPlacementButton.addEventListener('click', initializePlacementState);

startManualGameButton.addEventListener('click', async () => {
    const payload = {
        ships: placedShips.map(ship => ({
            Size: ship.Size,
            IsVertical: ship.IsVertical,
            Position: [{ X: ship.Position[0].X, Y: ship.Position[0].Y }],
        }))
    };
    try {
        const response = await fetch(`${API_URL}/newgame/manual`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });
        const result = await response.json();
        if (!response.ok) throw new Error(result.Message || 'Ошибка создания игры');
        placementContainer.style.display = 'none';
        mainGameContainer.style.display = 'flex';
        await updateGameView();
    } catch (error) {
        messageAreaEl.textContent = `Ошибка: ${error.message}`;
    }
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