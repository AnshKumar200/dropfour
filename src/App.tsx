import { Route, Routes } from 'react-router-dom'
import './App.css'
import IndexPage from './pages/IndexPage'
import GamePage from './pages/GamePage'

function App() {
    return (
        <Routes>
            <Route path='/' element={<IndexPage />}/>
            <Route path='/game/:id' element={<GamePage />}/>
        </Routes>
    )
}

export default App
