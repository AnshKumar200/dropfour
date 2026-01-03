let socket: WebSocket | null = null;

const WS_URL = import.meta.env.VITE_WS_URL ?? "localhost:7878";

export function connectWS(
    name: string,
    token: string | null,
    onMessage: (msg: any) => void
) {
    const params = new URLSearchParams()
    if (name) params.set("name", name)
    if (token) params.set("token", token)

    socket = new WebSocket(`ws://${WS_URL}/ws?${params}`)

    socket.onopen = () => {
        console.log("connected to ws")
    }

    socket.onmessage = (e) => {
        onMessage(JSON.parse(e.data))
    }
}

export function sendMove(column: number) {
    if (!socket) return;
    console.log("sent move:", column)
    socket.send(
        JSON.stringify({
            type: "move",
            data: { column }
        })
    )
}

export function joinQueue() {
    if(!socket) return;
    socket.send(
        JSON.stringify({
            type: "queue"
        })
    )
}
