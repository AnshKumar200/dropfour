import { sendMove } from "../ws";

type Props = {
    state: any;
}

export default function Game({ state }: Props) {
    console.log(state)
    return (
        <div>
            <div>turn: {state.Turn}</div>
            
            {state.Board.map((row: any[], r: number) => (
                <div key={r}>
                    {row.map((cell, c: number) => (
                        <button key={c} onClick={() => sendMove(c)}>
                            {cell === 0 ? "0" : cell === 1 ? "1" : "2"}
                        </button>
                    ))}
                </div>
            ))}
        </div>
    )
}
