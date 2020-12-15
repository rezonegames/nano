import Tetris from "./tetris";

export default class TetrisManager
{
    constructor()
    {
        this.instances = [];
    }

    createPlayer(element)
    {
        const tetris = new Tetris(element);

        this.instances.push(tetris);

        return tetris;
    }

    removePlayer(tetris)
    {
        this.instances = this.instances.filter(instance => instance !== tetris);
    }

    start()
    {
        this.instances.forEach((tetris)=>{
            tetris.run()
        })
    }
}
