export type Message = 
| { type: "start"; data: { player: number, token: string } }
| { type: "state"; data: GameState }
| { type: "resume"; data: GameState }
| { type: "end"; data: { winner: number } }

export type GameState = {
    Board: number[][];
    Turn: number;
    Players: any;
    Over: boolean;
}
