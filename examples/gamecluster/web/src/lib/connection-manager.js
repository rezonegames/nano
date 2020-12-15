export default class ConnectionManager
{
    constructor(conn, tetrisManager, uId)
    {
        this.uId = uId;
        this.conn = conn;
        this.peers = new Map;

        this.tetrisManager = tetrisManager;
        this.localTetris = this.tetrisManager.instances[0];

        this.watchEvents()
    }

    watchEvents()
    {
        const local = this.tetrisManager.instances[0];

        const player = local.player;
        ['pos', 'matrix', 'score'].forEach(key => {
            player.events.listen(key, () => {
                this.send('room.update', {
                    uId: this.uId,
                    type: 'state-update',
                    fragment: 'player',
                    state: [key, player[key]],
                });
            });
        });

        const arena = local.arena;
        ['matrix'].forEach(key => {
            arena.events.listen(key, () => {
                this.send('room.update', {
                    uId: this.uId,
                    type: 'state-update',
                    fragment: 'arena',
                    state: [key, arena[key]],
                });
            });
        });
    }

    updateManager(peers)
    {
        const clients = peers;
        clients.forEach(client => {
            if (!this.peers.has(client.uId)) {
                const tetris = this.tetrisManager.createPlayer(client.element);
                // tetris.unserialize(client.state);
                console.log("add peer:", client);
                this.peers.set(client.uId, tetris);
            }
        });

        [...this.peers.entries()].forEach(([uId, tetris]) => {
            if (!clients.some(client => client.uId === uId)) {
                this.tetrisManager.removePlayer(tetris);
                console.log("del peer:", uId);
                this.peers.delete(uId);
            }
        });

        console.log("after updateManager peers:", this.peers);
    }

    updatePeer(uId, fragment, [key, value])
    {
        if (!this.peers.has(uId)) {
            throw new Error('Client does not exist', uId);
            return
        }
        const tetris = this.peers.get(uId);
        tetris[fragment][key] = value;

        if (key === 'score') {
            tetris.updateScore(value);
        } else {
            tetris.draw();
        }
    }

    send(type, data)
    {
        // console.log('Sending message: type', type, 'data:', data);
        this.conn.request(type, data);
    }
}
