import { useState } from "react";

type Props = {
    onStart: (name: string) => void,
}

export default function Lobby({ onStart }: Props) {
    const [name, setName] = useState("");

    return (
        <div className="flex flex-col gap-2">
            <div>Lobby</div>
            <div className="flex gap-5">
                <input placeholder="Username" value={name} onChange={(e) => setName(e.target.value)} className="p-1 border-2 rounded-lg" />
                <button onClick={() => onStart(name)} className="bg-black text-white rounded-xl px-5">Connect</button>
            </div>
        </div>
    )
}
