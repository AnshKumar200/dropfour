import { useState } from "react";
import type { GameState, Message } from "../types";
import Lobby from "../components/Lobby";
import { connectWS } from "../ws";
import Game from "../components/Game";

export default function IndexPage() {
    const [inGame, setInGame] = useState(false);
    const [gameState, setGameState] = useState<GameState | null>(null);

    function handleStart(name: string) {
        const token = localStorage.getItem("token")

        connectWS(name, token, handleMessage);
    }

    function handleMessage(msg: Message) {
        if(msg.type === "start") {
            localStorage.setItem("token", msg.data.token)
            setInGame(true);
        }

        if(msg.type === "state" || msg.type === "resume") {
            console.log("recieved game state: ", msg.data)
            setGameState(msg.data)
            setInGame(true)
        }

        if(msg.type === "end") {
            alert(`winnder: ${msg.data.winner}`)
        }
    }

    if (!inGame) {
        return <Lobby onStart={handleStart} />
    }

    if(!gameState) {
        return <div>starting game...</div>
    }
    
    return <Game state={gameState} />
}
