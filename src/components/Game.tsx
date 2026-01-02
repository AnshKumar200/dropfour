import { sendMove } from "../ws";

type Props = {
    state: any;
}

export default function Game({ state }: Props) {
    return (
        <div>
            <div>turn: {state.Turn}</div>
            
            {state.Board.map((row: any[], r: number) => (
                <div key={r} className="flex">
                    {row.map((cell, c: number) => (
                        <div key={c} onClick={() => sendMove(c)} className={`${cell === 0 ? "bg-white" : cell === 1 ? "bg-blue-500" : "bg-red-500"} size-5 border-2 rounded-full`}>
                        </div>
                    ))}
                </div>
            ))}
        </div>
    )
}
