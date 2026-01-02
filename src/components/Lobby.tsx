import { useState } from "react";

type Props = {
    onStart: (name: string) => void,
}

export default function Lobby({ onStart }: Props) {
    const [name, setName] = useState("");

    return (
        <div>
            <div>lobby</div>
            <input placeholder="username" value={name} onChange={(e) => setName(e.target.value)} />
            <button onClick={() => onStart(name)}>play</button>
        </div>
    )
}
