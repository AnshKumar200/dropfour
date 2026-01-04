import { useEffect, useState } from "react";
import { gameJoinQueue, sendMove } from "../ws";

type Props = {
    state: any;
}

const TURN_TIME = 30;

export default function Game({ state }: Props) {
    const [timeLeft, setTimeLeft] = useState(TURN_TIME);
    const [inQueue, setInQueue] = useState(false)

    useEffect(() => {
        if (!state.Over) {
            setInQueue(false)
        }
    }, [state.Over])

    function handleQueue() {
        if (inQueue) return
        setInQueue(true)
        gameJoinQueue()
    }

    useEffect(() => {
        if (!state.LastMoveTime || state.Over) return;

        const updateTimer = () => {
            const now = Date.now();
            const elapsed = Math.floor((now - state.LastMoveTime) / 1000);
            const rem = Math.max(TURN_TIME - elapsed, 0);
            setTimeLeft(rem);
        }

        updateTimer()
        const interval = setInterval(updateTimer, 1000)

        return () => clearInterval(interval)
    }, [state.LastMoveTime, state.Over])


    return (
        <div className="text-nowrap">
            <div>Player 1: {state.Players[0].Name} V/S Player 2: {state.Players[1].Name}</div>
            <div>Turn: Player {state.Turn}</div>
            <div>Time Left: {timeLeft}s</div>

            {state.Board.map((row: any[], r: number) => (
                <div key={r} className="flex">
                    {row.map((cell, c: number) => (
                        <div key={c} onClick={() => sendMove(c)} className={`${cell === 0 ? "bg-white" : cell === 1 ? "bg-blue-500" : "bg-red-500"} size-8 border-2 rounded-full`}>
                        </div>
                    ))}
                </div>
            ))}
            {state.Over && <div>
                {state.Winner !== 0 ? (
                    <div>Winner is: Player {state.Winner} - {state.Players[state.Winner - 1].Name}</div>
                ) : (
                    <div>Draw</div>
                )}
                <button onClick={handleQueue} className="px-5 py-2 bg-black text-white rounded-xl">Play Again</button>
                {inQueue && <div>you are in queue!</div>}
            </div>}
        </div>
    )
}
