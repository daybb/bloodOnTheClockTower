import React, {useEffect, useRef, useState} from "react"
import {Button, Card, notification, Select, Switch} from "antd"
import "./ClockTower.css"
import config from "../../config/config"
import {useParams} from "react-router-dom"
import {sleep} from "../../utils/time"
import {remove} from "../../utils/array"

const Context = React.createContext({
    name: "Default",
})
let socket
let castLock = false // 入夜后给施放技能的时间，后端没有，只在前端限制，因为只限制主机

// 角色技能名称映射
const skillNameMap = {
    "Assassin": "刺杀",
    "Gambler": "猜角色",
    "Professor": "复活",
    "Fool": "混淆",
    "Shabaloth": "吞食"
}

function ClockTower() {
    let { roomId } = useParams()
    // 加载游戏
    const [game, setGame] = useState(null)
    useEffect(() => {
        establishConn()
        return () => {
            if (socket) {
                socket.close()
            }
        }
    }, [roomId])
    const establishConn = () => {
        // 获取game 长连接
        socket = new WebSocket(`${config.beBaseUrl}/game/${roomId}/${sessionStorage.getItem("PlayerID")}`)
        socket.onopen = function() {
            loadGame()
        }
        socket.onmessage = function(event) {
            // TODO 内测用，记得关闭
            // console.log("Received message from server:", JSON.parse(event.data))
            setGame(JSON.parse(event.data))
        }
        socket.onerror = function() {
            establishConn() // 断线重连
        }
    }
    const loadGame = () => {
        let req = JSON.stringify({action: "load_game", targets: []})
        socket.send(req)
    }

    //*************** 用户展示log相关***************//
    useEffect(() => {
        loadPersonalLog()
        // updateSeatTag()
        // updateSeatDead()
    }, [game])
    const loadPersonalLog = async () => {
        if (game) {
            for (let i = 0; i < game.players.length; i++) {
                if (game.players[i].id === sessionStorage.getItem("PlayerID")) {
                    replaceLog(game.players[i].log || "暂无游戏日志", ...wordClassPairs)
                    break
                }
            }
        }
    }
    const updateText = (text, word, className) => {
        if (typeof word === "string") {
            let regex = new RegExp(word, "g")
            return text.replace(regex, `<span class="${className}">${word}</span>`)
        }
        return text.replace(word, match => `<span class="${className}">${match}</span>`)
    }
    const replaceLog = (text, ...wordClassPairs) => {
        if (!text) return
        
        let replacedText = updateText(text, ...wordClassPairs[0])
        if (wordClassPairs.length > 1) {
            for (let i = 1; i < wordClassPairs.length; i++) {
                replacedText = updateText(replacedText, ...wordClassPairs[i])
            }
        }
        let removedNTextArr = replacedText.split("\n").map(item => {
            return `<span>${item}</span>`
        })
        let resultText = removedNTextArr.join("")
        if (document.getElementById("LOG")) {
            document.getElementById("LOG").innerHTML = `${resultText}`
        }
    }
    let wordClassPairs = [
        [/(?<=第).*?(?=天)|平安夜/g, "highlight highlight-number"], // 数字
        [/\[([^\]]+)]/g, "highlight highlight-player"], // 玩家名字
        [/\{[^}]+}/g, "highlight highlight-skill-result"], // 技能结果关键字
        [/(投毒|卜算|认主|守护|杀害|枪毙|反弹|反向通灵|注册|投给)/g, "highlight highlight-skill"], // 技能关键字
        [/(死亡|处决结果)/g, "highlight highlight-severe"], // 重大事件关键字
        [/(提名)/g, "highlight highlight-nominate"], // 提名
        [/(投票)/g, "highlight highlight-vote"], // 投票
    ]

    // 玩家数据
    const [players, setPlayers] = useState([])
    // 当前玩家
    const [currentPlayer, setCurrentPlayer] = useState(null)
    // 当前回合
    const [currentRound, setCurrentRound] = useState(1)
    // 技能使用记录
    const [skillUsed, setSkillUsed] = useState({})
    // 通知API
    const [api, contextHolder] = notification.useNotification()
    // 是否显示自己的身份
    const [showMyRole, setShowMyRole] = useState(true)
    //显示玩家信息
    const findPlayer = () => {
        if (game !== null) {
            for (let i = 0; i < game.players.length; i++) {
                if (game.players[i].id === sessionStorage.getItem("PlayerID")) {
                    return {
                        name: game.players[i].name,
                        character: game.players[i].baseCharacter.character_name,
                        characterType: game.players[i].characterType,
                        positionId: game.players[i].positionId
                    }
                }
            }
        }
        return {}
    }
    //*************角色使用技能****************//
    const getMe = (game) => {
        let me
        for (let i = 0; i < game.players.length; i++) {
            if (game.players[i].id === sessionStorage.getItem("PlayerID")) {
                me = game.players[i]
                break
            }
        }
        return me
    }
    // 点击玩家名字，选中玩家，保存被选中的玩家ID
    const [selectedPlayers, setSelectedPlayers] = useState([])
    const selectPlayer = (event) => {
        event.preventDefault()
        // 当前玩家没死，就可以选择玩家，提名或施放技能
        let me = getMe(game)
        console.log("me",me.baseCharacter.is_dead)
        if (!me.baseCharacter.is_dead) {
            let selectedPlayersCopy = selectedPlayers.slice()
            if (event.target.classList.contains("seat-selected")) {
                event.target.classList.remove("seat-selected")
                remove(selectedPlayersCopy, event.target.id)
                setSelectedPlayers(selectedPlayersCopy)
            } else {
                console.log("add now")
                event.target.classList.add("seat-selected")
                selectedPlayersCopy.push(event.target.id)
                setSelectedPlayers(selectedPlayersCopy)
            }
        }
    }
    const [selectedCharacter, setSelectedCharacter] = useState(null)
    const [isCastModalOpen, setIsCastModalOpen] = useState(false)
    const [castModalContent, setCastModalContent] = useState("抱歉，您无法发动技能")
    const showCastModal = () => {
        setIsCastModalOpen(true)
        setCastModalContent("确定发动技能吗") // 将modal的内容重新初始化，防止错乱
        let me = getMe(game)
        setCastModalContent(genCastModalContent(me))
    }
    // 产生技能施放Modal的内容
    const genCastModalContent = (me) => {
        const [isHidden] = useState(false)
        const inputRef = useRef(null)
        useEffect(() => {
            if (isHidden && inputRef.current) {
                inputRef.current.blur() // 隐藏时让输入框失去焦点
            }
        }, [isHidden])
        if (game.state.stage === 0) {
            return "本局未开始，不能发动技能"
        }
        if (me.state.dead) {
            return "您已死亡"
        }
        if (game.state.votingStep && me.character !== "杀手") {
            return "投票阶段不能发动技能"
        }
        // Gambler 特殊处理
        // Gambler 特殊处理
        if (me.baseCharacter.character_name === "Gambler") {
            if (hasGuessed) {
                return "您本局游戏已经使用过猜测能力，每局只能猜测一次。"
            }

            return (
                <div
                    style={{ display: isHidden ? "none" : "block" }}
                    aria-hidden={isHidden}
                >
                    <p>请选择您要猜测的角色：</p>
                    <Select
                        ref={inputRef}
                        style={{ width: "100%" }}
                        value={selectedCharacter}
                        onChange={(value) => setSelectedCharacter(value)}
                        options={characterOptions.map((character) => ({
                            value: character,
                            label: character
                        }))}
                        // 当元素隐藏时自动失去焦点
                        onBlur={() => {
                            if(isHidden && inputRef.current) {
                                inputRef.current.blur()
                            }
                        }}
                    />
                </div>
            )
        }

        let selectedPlayersObj = []
        let content = "是否要对玩家 "
        for (let i = 0; i < selectedPlayers.length; i++) {
            for (let j = 0; j < game.players.length; j++) {
                if (selectedPlayers[i] === game.players[j].id) {
                    content += "<" + game.players[j].name + "> "
                    selectedPlayersObj.push(game.players[j])
                    break
                }
            }
        }
        switch (me.baseCharacter.character_name) {
            case "Gambler":
                if (selectedPlayers.length === 1) {
                    content += `猜测为 ${selectedCharacter || "未选择"} 角色吗？`
                    break
                }
                return "您只能选1个人进行猜测"
            case "Shabaloth":
                if (selectedPlayers.length === 2) {
                    content += "占卜，看看有没有恶魔吗？"
                    break
                }
                return "您只能选2个人占卜"
            case "Assassin":
                if (!game.state.night) {
                    return "白天不能开枪"
                }
                // if (game.executed) {
                //     return "已发生处决，不能开枪"
                // }
                if (selectedPlayers.length === 1) {
                    content += "实行刺杀吗？"
                    break
                }
                return "您只能选1个人刺杀"
        }
        setCastToPlayersId(selectedPlayers)
        return content
    }
    // 使用技能
    const handleCastOk = () => {
        showCastModal()
        console.log("cast ok")
        // setIsCastModalOpen(false)
        // 检查是否已经在本回合使用过技能
        if (skillUsed[currentPlayer.id]) {
            console.log("cast fail")
            openNotification("技能使用失败", "您在本回合已经使用过技能")
            return
        }
        // 后端判断 发动技能的条件是，取决于身份，drunk，白天黑夜，还有没有技能；前端随便发动，后端判断成不成功
        let me = getMe(game)
        console.log("me name",me.baseCharacter.character_name)
        switch (me.baseCharacter.character_name){
        //沙巴螺丝攻击判定
        case "Shabaloth":
            if (selectedPlayers.length !== 2) {
                openNotification("技能使用失败", "请选择两名目标")
                return
            }
            {
                let req = JSON.stringify({
                    action: "cast",
                    targets: selectedPlayers
                })
                socket.send(req)
                setSelectedCharacter(null) // 重置选择
            }
            break
        //赌徒攻击判定
        case "Gambler":
            if (selectedCharacter && selectedPlayers.length === 1){
                //弹出选择框选择猜测的角色
                let req = JSON.stringify({
                    action: "cast",
                    targets: selectedPlayers,
                    extra: selectedCharacter // 传递选中的角色名称
                })
                socket.send(req)
                setSelectedCharacter(null) // 重置选择
                break
            }
            break
        default:
            api.info({
                message: me.baseCharacter.character_name+"技能使用",
                description: "您本局游戏已经使用过猜测能力，每局只能猜测一次。",
                placement: "topRight",
            })
        }
        // 标记技能已使用
        setSkillUsed({
            ...skillUsed,
            [currentPlayer.id]: true
        })
        setSelectedPlayers([])
    }

    //首次入夜
    const emitCheckOutFirstNight = () => {
        let req = JSON.stringify({action: "first_night", targets: []})
        socket.send(req) // 会在后端更新stage、night
    }
    // 初始化玩家数据
    useEffect(() => {
        if (game && game.players && game.players.length > 0) {
            // 从game.players获取真实玩家数据
            const realPlayers = game.players.map(player => ({
                id: player.id,
                name: player.name,
                role: player.baseCharacter.character_name,
                isDead: player.isDead || false,
                isNominated: player.isNominated || false,
                characterType: player.characterType,
                positionId: player.positionId
            }))

            setPlayers(realPlayers)
            
            // 设置当前玩家
            const currentPlayerID = sessionStorage.getItem("PlayerID")
            const myPlayer = realPlayers.find(p => p.id === currentPlayerID)
            if (myPlayer) {
                setCurrentPlayer(myPlayer)
            }
        }
    }, [game])

    // 重置技能使用状态（每回合开始时）
    useEffect(() => {
        setSkillUsed({})
    }, [currentRound])

    // 提名玩家
    const nominate = () => {
        if (selectedPlayers.length !== 1) {
            openNotification("提名失败", "请选择一名玩家进行提名")
            return
        }

        const nominatedPlayerId = selectedPlayers[0]
        setPlayers(players.map(player =>
            player.id === nominatedPlayerId
                ? {...player, isNominated: true}
                : player
        ))

        openNotification("提名成功", `已提名玩家: ${players.find(p => p.id === nominatedPlayerId).name}`)
        setSelectedPlayers([])
    }

    // 投票
    const vote = () => {
        const nominatedPlayers = players.filter(player => player.isNominated)

        if (nominatedPlayers.length === 0) {
            openNotification("投票失败", "当前没有被提名的玩家")
            return
        }

        // 模拟投票结果
        openNotification("投票成功", "您的投票已记录")
    }

    const checkReadyToToggleNight = () => {
        //第一个白天可以不需要发动技能
        // 死亡或者已放过技能都是ready
        let ready = true
        for (let i = 0; i < game.players.length; i++) {
            ready = ready && (game.players[i].state.casted || game.players[i].state.dead)
        }
        // 所有有技能的操作完，没技能的点完验证码，时间等待结束，不在投票阶段，则切换日夜，切换后首先结算前一阶段
        return !game.state.votingStep
            && !castLock
            && ready // TODO 测试时，可注释
            || game.state.stage === 0 || game.state.stage ===1
    }
    //发送checkout_night到后端
    const emitToggleNight = () => {
        let req = JSON.stringify({action: "checkout_night", targets: []})
        socket.send(req) // 会在后端更新stage、night
    }
    // 游戏过程
    const gameProcess = async () => {
        castLock = true
        // 发送日夜切换指令到后端，后端重置状态
        await emitToggleNight()
        // 防抖
        await sleep(2000)
        castLock = false
    }
    // 结束日/夜
    const toggleDayNight = () => {
        //结束日
        if(!game || !game.state) return

        if(!game.state.night){
            //进入夜晚
            let req = JSON.stringify({action: "direct_night", targets: []})
            socket.send(req) // 会在后端更新stage、night
        }else {
            //进入白天
            if (checkReadyToToggleNight()) {
                // 锁定与结算过程
                gameProcess(game.state.stage+1)
                setCurrentRound(currentRound + 1)
                // 重置提名状态
                setPlayers(players.map(player => ({...player, isNominated: false})))
                setSelectedPlayers([])
                openNotification("回合结束", `第${currentRound}回合已结束，开始第${currentRound + 1}回合`)
            } else {
                openNotification("切换失败", "有玩家尚未完成操作")
                return
            }
        }
        openNotification("时间切换", !game.state.night ? "已切换到夜晚" : "已切换到白天")

    }
    // 结束投票
    const emitEndVoting = () => {
        let req = JSON.stringify({action: "end_voting", targets: []})
        socket.send(req) // 会在后端更新stage、night
    }
    const openEndVotingNotification = (placement) => {
        api.info({
            message: "非法点击",
            description: <Context.Consumer>{() => "不好意思, 当前不在投票环节，无法结束投票环节!"}</Context.Consumer>,
            placement,
        })
    }
    //除第一天之外，处决并入夜。无需处决也点这个按钮入夜
    const emitExecute = () => {
        if (game && game.state.votingStep) {
            emitEndVoting()
        } else {
            openEndVotingNotification("topRight")
        }

        let req = JSON.stringify({action: "execute", targets: []})
        socket.send(req) // 会在后端更新stage、night
    }
    // 结束回合
    // const endRound = () => {
    //     setCurrentRound(currentRound + 1)
    //     // 重置提名状态
    //     setPlayers(players.map(player => ({...player, isNominated: false})))
    //     setSelectedPlayers([])
    //     openNotification("回合结束", `第${currentRound}回合已结束，开始第${currentRound + 1}回合`)
    // }

    // 通知函数
    const openNotification = (title, description) => {
        api.info({
            message: title,
            description,
            placement: "topRight"
        })
    }

    // 获取当前玩家的技能名称
    const getSkillName = () => {
        if (!currentPlayer) return "无技能"
        return skillNameMap[currentPlayer.role] || "无技能"
    }

    // 计算玩家位置的角度
    const getPlayerPositions = () => {
        return players.map((player, index) => {
            const angle = (index * (360 / players.length)) * (Math.PI / 180)
            const radius = 40 // 圆的半径（百分比）
            const x = 50 + radius * Math.cos(angle)
            const y = 50 + radius * Math.sin(angle)
            return {...player, position: {x, y}}
        })
    }

    const positionedPlayers = getPlayerPositions()

    return (
        <div className="clock-tower"
            style={{backgroundColor: game && game.state && !game.state.night ? "#f0f2f5" : "#001529"}}>
            {contextHolder}

            {/* 游戏信息 */}
            <div className="game-info">
                <h1>血染钟楼</h1>
                <div className="round-info">
                    <span>第 {currentRound} 回合</span>
                    <span
                        className="day-night-indicator">{game && game.state ? (!game.state.night ? "白天" : "夜晚") : "等待中"}</span>
                </div>
            </div>
            <div className="layout east" id="LOG">
                <span>未开始，尚未入第一夜</span>
            </div>
            {/* 左侧玩家信息面板 */}
            <div className="player-info-panel">
                <Card title="玩家信息" bordered={false} className="player-info-card">
                    {currentPlayer && (
                        <>
                            <div className="info-item">
                                <span className="info-label">昵称：</span>
                                <span className="info-value">{currentPlayer.name}</span>
                            </div>
                            <div className="info-item">
                                <span className="info-label">位置：</span>
                                <span
                                    className="info-value">{players.findIndex(p => p.id === currentPlayer.id) + 1}号位</span>
                            </div>
                            <div className="info-item">
                                <span className="info-label">身份：</span>
                                <span className="info-value">
                                    {showMyRole ? currentPlayer.role : "******"}
                                </span>
                            </div>
                            <div className="info-item role-toggle">
                                <span className="info-label">显示身份：</span>
                                <Switch checked={showMyRole} onChange={setShowMyRole}/>
                            </div>
                            <div className="info-item">
                                <span className="info-label">状态：</span>
                                <span className="info-value">
                                    {currentPlayer.isDead ? "已死亡" : "存活中"}
                                    {currentPlayer.isNominated ? "，被提名" : ""}
                                </span>
                            </div>
                        </>
                    )}
                </Card>
            </div>

            {findPlayer().positionId === 1 && game.state.stage === 0 && !game.state.night &&
                <Button className="btn small-btn" onClick={emitCheckOutFirstNight}>进入第一夜</Button>}
            {findPlayer().positionId === 1 && !game.state.night && game.executed &&
                <Button className="btn small-btn" onClick={emitExecute}>投票结束，处决并入夜</Button>}
            {/*有处决的时候是处决，没有的时候是入夜*/}
            {findPlayer().positionId === 1 && <Button type="primary" onClick={toggleDayNight}>
                {game && game.state ? (!game.state.night && !game.executed ? "直接入夜" : "天亮") : "切换时间"}
            </Button>}
            {/* 玩家圆形布局 */}
            <div className="players-circle">
                {positionedPlayers.map((player) => (
                    <div
                        key={player.id}
                        className={`player-seat ${selectedPlayers.includes(player.id) ? "selected" : ""} ${player.isDead ? "dead" : ""} ${player.isNominated ? "nominated" : ""}`}
                        style={{
                            left: `${player.position.x}%`,
                            top: `${player.position.y}%`,
                        }}
                        onClick={selectPlayer}
                    >
                        <div className="player-name">{player.name}</div>
                        {player.isDead && <div className="death-status">已死亡</div>}
                        {player.isNominated && <div className="nomination-badge">被提名</div>}
                    </div>
                ))}

                {/* 操作按钮 - 放在圆形布局中间 */}
                <div className="center-action-buttons">
                    <Button type="primary" onClick={nominate}>提名</Button>
                    <Button type="primary" onClick={vote}>投票</Button>
                    {/*<Button type="primary" onClick={endRound}>结束回合</Button>*/}
                    <Button
                        type="primary"
                        open={isCastModalOpen} onOk={handleCastOk}
                        disabled={skillUsed[currentPlayer?.id]}
                        className={skillUsed[currentPlayer?.id] ? "skill-used" : ""}
                    >
                        {getSkillName()}
                    </Button>
                </div>
            </div>
        </div>
    )
}

export default ClockTower
