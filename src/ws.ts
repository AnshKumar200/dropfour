let socket: WebSocket | null = null;
let connected = false;

const WS_URL = import.meta.env.VITE_WS_URL ?? "localhost:7878";

export function connectWS(
    name: string,
    token: string | null,
    onMessage: (msg: any) => void,
    onClose?: () => void,
) {
    if (socket) {
        return;
    }

    const params = new URLSearchParams()
    if (name) params.set("name", name)
    if (token) params.set("token", token)

    socket = new WebSocket(`wss://${WS_URL}/ws?${params}`)

    socket.onopen = () => {
        connected = true;
        console.log("connected to ws")

        getLeaderboard()
    }

    socket.onmessage = (e) => {
        onMessage(JSON.parse(e.data))
    }

    socket.onerror = (e) => {
        console.error("websocket error: ", e)
    }

    socket.onclose = () => {
        console.log("websocket closed")
        connected = false
        socket = null
        onClose?.()
    }
}

export function sendMove(column: number) {
    if (!connected || !socket) return;
    console.log("sent move:", column)
    socket.send(
        JSON.stringify({
            type: "move",
            data: { column }
        })
    )
}

export function joinQueue() {
    if (!socket) return;
    socket.send(
        JSON.stringify({
            type: "queue"
        })
    )
}

export function gameJoinQueue() {
    if (!connected || !socket) return;
    socket.send(
        JSON.stringify({
            type: "game_queue"
        })
    )
}

export function getLeaderboard() {
    if (!connected || !socket) return;
    socket.send(
        JSON.stringify({
            type: "leaderboard"
        })
    )
}
