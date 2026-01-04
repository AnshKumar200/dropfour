import { useState } from "react";
import { type GamesData, type GameState, type LeaderboardData, type Message } from "./types";
import Lobby from "./components/Lobby";
import { connectWS, getData, joinQueue } from "./ws";
import Game from "./components/Game";

export default function IndexPage() {
    const [connected, setConnected] = useState(false);
    const [inGame, setInGame] = useState(false);
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [leaderboard, setLeaderboard] = useState<LeaderboardData[] | null>(null);
    const [inQueue, setInQueue] = useState(false);
    const [gamesData, setGamesData] = useState<GamesData[] | null>(null);

    function handleStart(name: string) {
        const token = localStorage.getItem("token") ?? "";
        connectWS(name, token, handleMessage);
        setConnected(true)
    }

    function handlePlay() {
        joinQueue()
        setInQueue(true)
    }

    function handleMessage(msg: Message) {
        switch (msg.type) {
            case "token":
                localStorage.setItem("token", msg.data)
                break;
            case "start":
                setInGame(true);
                break;
            case "state":
            case "resume":
                setGameState(msg.data)
                setInGame(true)
                setInQueue(false)
                break;
            case "leaderboard":
                setLeaderboard(msg.data)
                break;
            case "games":
                setGamesData(msg.data)
                break;
            case "end":
                alert(`winner: ${msg.data.winner}`)
                break;
        }
    }

    function handleUpdate() {
        getData()
    }

    return (
        <div className="flex gap-10">
            <div>
                {!connected && <Lobby onStart={handleStart} />}
                {connected && !inGame && (
                    <div>
                        <div>Connected!</div>
                        <button onClick={handlePlay} className="bg-black text-white px-5 py-2 rounded-xl">Play</button>
                        <div>{inQueue ? "you are in the queue!" : ""}</div>
                    </div>
                )}
                {inGame && !gameState && <div>starting game...</div>}
                {inGame && gameState && <Game state={gameState} />}
            </div>
            <div>
                <button onClick={handleUpdate} className="px-5 py-2 bg-black text-white rounded-xl">Update Data</button>
                <div className="flex gap-10">
                    <div>
                        <div>Leaderboard</div>
                        <div>
                            {leaderboard && leaderboard.map((player: LeaderboardData, key: number) => (
                                <div key={key} className="flex gap-10 text-nowrap">
                                    <div>Name: {player.Name}</div>
                                    <div>Wins: {player.Wins}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                    <div>
                        <div>Games</div>
                        <div className="flex gap-10 flex-wrap">
                            {gamesData && gamesData.map((game, key) => (
                                <div key={key}>
                                    <div>Player 1: {game.Player1}</div>
                                    <div>Player 2: {game.Player2}</div>
                                    <div>Winner: Player {game.Winner}</div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}
