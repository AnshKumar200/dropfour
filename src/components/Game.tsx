import { joinQueue, sendMove } from "../ws";

type Props = {
    state: any;
}

export default function Game({ state }: Props) {
    return (
        <div>
            <div>Player 1: {state.Players[0].Name} V/S Player 2: {state.Players[1].Name}</div>
            <div>Turn: {state.Turn}</div>

            {state.Board.map((row: any[], r: number) => (
                <div key={r} className="flex">
                    {row.map((cell, c: number) => (
                        <div key={c} onClick={() => sendMove(c)} className={`${cell === 0 ? "bg-white" : cell === 1 ? "bg-blue-500" : "bg-red-500"} size-5 border-2 rounded-full`}>
                        </div>
                    ))}
                </div>
            ))}
            {state.Over && (
                <div>
                    <div>Winner is: Player {state.Winner} - {state.Players[state.Winner - 1].Name}</div>
                    <button onClick={joinQueue}>Play Again</button>
                </div>
            )}
        </div>
    )
}
