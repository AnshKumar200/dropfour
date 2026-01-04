import { useState } from "react";
import type { GameState, LeaderboardData, Message } from "./types";
import Lobby from "./components/Lobby";
import { connectWS, joinQueue } from "./ws";
import Game from "./components/Game";

export default function IndexPage() {
    const [connected, setConnected] = useState(false);
    const [inGame, setInGame] = useState(false);
    const [gameState, setGameState] = useState<GameState | null>(null);
    const [leaderboard, setLeaderboard] = useState<LeaderboardData[] | null>(null);
    const [inQueue, setInQueue] = useState(false);

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
                console.log("got leader: ", msg.data)
                setLeaderboard(msg.data)
                break;
            case "end":
                alert(`winner: ${msg.data.winner}`)
                break;
        }
    }

    return (
        <div className="flex gap-50">
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
                <div>Leaderboard</div>
                <div>
                    {leaderboard && leaderboard.map((player: LeaderboardData, key: number) => (
                        <div key={key} className="flex gap-10">
                            <div>Name: {player.Name}</div>
                            <div>Wins: {player.Wins}</div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    )
}
