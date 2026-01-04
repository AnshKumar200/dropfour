export type Message = 
| { type: "token"; data: string }
| { type: "start"; }
| { type: "state"; data: GameState }
| { type: "resume"; data: GameState }
| { type: "leaderboard"; data: LeaderboardData[] }
| { type: "games"; data: GamesData[] }
| { type: "end"; data: { winner: number } }

export type GameState = {
    Board: number[][];
    Turn: number;
    Players: any;
    Winner: number;
    Over: boolean;
    LastMoveTime: number;
}

export type LeaderboardData = {
    Name: string;
    Wins: number;
}

export type GamesData = {
    Player1: string;
    Player2: string;
    Winner: number;
}
