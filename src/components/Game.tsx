import { useEffect, useState } from "react";
import { joinQueue, sendMove } from "../ws";

type Props = {
    state: any;
}

const TURN_TIME = 30;

export default function Game({ state }: Props) {
    const [timeLeft, setTimeLeft] = useState(TURN_TIME);

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
        <div>
            <div>Player 1: {state.Players[0].Name} V/S Player 2: {state.Players[1].Name}</div>
            <div>Turn: Player {state.Turn}</div>
            <div>Time Left: {timeLeft}s</div>

            {state.Board.map((row: any[], r: number) => (
                <div key={r} className="flex">
                    {row.map((cell, c: number) => (
                        <div key={c} onClick={() => sendMove(c)} className={`${cell === 0 ? "bg-white" : cell === 1 ? "bg-blue-500" : "bg-red-500"} size-5 border-2 rounded-full`}>
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
                <button onClick={joinQueue}>Play Again</button>
            </div>}
        </div>
    )
}
