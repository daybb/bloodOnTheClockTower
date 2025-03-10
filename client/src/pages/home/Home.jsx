/* eslint-disable */
import React, {useEffect, useMemo, useState} from "react"
import "./Home.css"
import {HOME} from "../../config/cstModule"
import {Avatar, Button, Drawer, Flex, Input, List, notification, Space} from "antd"
import {genShortUUID} from "../../utils/uuid"
import {useNavigate} from "react-router-dom"
import {blood} from "../../utils/blood/blood"
import {sleep} from "../../utils/time"
import config from "../../config/config"
import bat from "../../assets/video/bat.gif"

if (sessionStorage.getItem("PlayerID") === null) {
    sessionStorage.setItem("PlayerID", genShortUUID())
}
if (sessionStorage.getItem("PlayerName") === null) {
    sessionStorage.setItem("PlayerName", "好人1号-" + genShortUUID().slice(-6))
}

const Context = React.createContext({
    name: "Default",
})

let socketHome


function Home() {
    const [roomList, setRoomList] = useState([])

    useEffect(() => {
        establishConn()
    }, [])
    const establishConn = () => {
        socketHome = new WebSocket(`${config.beBaseUrl}/home`)
        socketHome.onopen = function () {
            loadRoomList()
        }
        socketHome.onmessage = function (event) {
            // console.log("Received message from server:", JSON.parse(event.data))
            setRoomList(JSON.parse(event.data))
        }
        socketHome.onerror = function (error) {
            console.error("WebSocket error:", error)
        }
    }
    const loadRoomList = () => {
        let data = {
            action: "list_rooms",
            payload: sessionStorage.getItem("PlayerID"),
        }
        socketHome.send(JSON.stringify(data))
    }
    const navigate = useNavigate()
    const jump = (roomId,playerId) => {
        navigate(`/room/${roomId}/${playerId}`, {
            replace: true,
            state: {roomId}
        })
    }

    const [open, setOpen] = useState(false)

    const onClose = () => {
        setOpen(false)
    }

    const [open1, setOpen1] = useState(false)
    const showDrawer1 = room => {
        console.log("name is: ", sessionStorage.getItem("PlayerName"))
        // console.log("room list is: ", roomList)
        setOpen1(true)
        setRoomId(roomId)
        setRoomName(roomName)
        setRoomSelected(room)
        setPlayerName(sessionStorage.getItem("PlayerName"))
    }
    const onClose1 = () => {
        setOpen1(false)
    }

    const [roomSelected, setRoomSelected] = useState(null)

    //create room 请求后端建立新房间，刷新当前页面
    const createRoom = () => {
        onClose()
        let player = {}
        player.id = playerName
        player.name = playerName
        player.waiting = true
        let roomInfo = {
            id: roomId,
            name: roomName,
            host: "1",
            // players: [player],
        }
        let data = {
            action: "create_room",
            payload: roomInfo,
        }
        socketHome.send(JSON.stringify(data))
        window.location.reload(); // 刷新页面
        // jump(roomId)
    }
    const [seatPositionId, setSeatPosition] = useState(1)
    const [playerName, setPlayerName] = useState("鸡掰")
    const [roomId, setRoomId] = useState("2025")
    const [roomName, setRoomName] = useState("摸摸鱼")
    const [error, setError] = useState('');
    // console.log(sessionStorage.getItem("PlayerID"))
    const handsetSeatPosition = (event) => {
        const value = event.target.value;
        if (/^-?\d*$/.test(value)) { // 检查是否为整数
            setSeatPosition(value === '' ? '' : parseInt(value, 10)); // 转换为整数
            setError('');
        } else {
            setError('请输入有效的整数');
        }
        // setSeatPosition(event.target.value)
    }

    const handlePlayerNameChange = (event) => {
        setPlayerName(event.target.value)
        sessionStorage.setItem("PlayerName", event.target.value)
    }

    const joinRoom = () => {
        // 如果房间status是游戏中，则无法加入
        // if (roomSelected && (roomSelected.status === "游戏中" || roomSelected.status === "复盘中")) {
        //     openGamingNotification("topRight")
        //     return
        // }
        // console.log(roomSelected)
        // 人数大于等于15，则无法加入房间
        if (roomSelected && roomSelected.players.length >= 15) {
            openPlayerNum2Notification("topRight")
            return
        }
        onClose1()
        let playerInfo = {
            id: playerName,
            name: playerName,
            waiting: true,
            positionId: seatPositionId,
        }
        console.log("position: ", seatPositionId)
        let roomInfo = {}
        let data = {
            action: "join_room",
            payload: {
                room: roomInfo,
                player: playerInfo,
            },
        }
        socketHome.send(JSON.stringify(data))
        socketHome.onmessage = function (event) {
            // setRoomList(JSON.parse(event.data))
        }
        sessionStorage.setItem("PlayerID", playerInfo.id)
        jump(roomId,playerInfo.id)
    }
    const [api, contextHolder] = notification.useNotification()
    const openPlayerNum2Notification = (placement) => {
        api.info({
            message: "操作无效",
            description: <Context.Consumer>{() => "人数多于十五人，无法加入房间!"}</Context.Consumer>,
            placement,
        })
    }
    const openGamingNotification = (placement) => {
        api.info({
            message: "操作无效",
            description: <Context.Consumer>{() => "游戏已开始或还在复盘，无法加入房间!"}</Context.Consumer>,
            placement,
        })
    }
    const openRoomNameNotification = (placement) => {
        api.info({
            message: "信息确实",
            description: <Context.Consumer>{() => "房间名称不能为空!"}</Context.Consumer>,
            placement,
        })
    }
    const contextValue = useMemo(
        () => ({
            name: "",
        }),
        [],
    )

    // 动画 蝙蝠 流血
    useEffect(() => {
        blood()
        hideGif("Bat-gif", 3000)
    }, [])
    const hideGif = async (id, ms) => {
        let gif = document.getElementById(id)
        await sleep(ms)
        if (gif) {
            gif.classList.add("hidden")
        }
    }

    return (
        <div id={HOME.KEY} className={HOME.KEY}>
            <div id="Title-wrap">
                <div id="Title">血染钟楼</div>
            </div>
            <svg className="svg">
                <filter id="noise">
                    <feTurbulence baseFrequency="0.07" type="fractalNoise" result="turbNoise"></feTurbulence>
                    <feDisplacementMap in="SourceGraphic" in2="turbNoise" xChannelSelector="G" yChannelSelector="B"
                                       scale="6" result="disp"></feDisplacementMap>
                </filter>
            </svg>
            {window.innerWidth >= 958 ?
                <img id="Bat-gif" src={bat} alt="Bat GIF"/>
                :
                <></>
            }
            <Flex className="layout" wrap="wrap">
                <Button className="btn-main" onClick={createRoom}>创建房间</Button>
            </Flex>
            <div className="layout">
                {roomList ?
                    <List style={{width: "100%"}}
                          itemLayout="horizontal"
                          dataSource={roomList}
                          renderItem={(item, index) => (
                              <List.Item className="list-item">
                                  <List.Item.Meta
                                      avatar={<Avatar src={`https://api.dicebear.com/7.x/miniavs/svg?seed=${index}`}/>}
                                      title={
                                          <span>{item.name}</span>
                                      }
                                      description={
                                          <Flex horizontal="true" gap="middle" justify="space-between" align="center"
                                                wrap="wrap">
                                              <span>当前人数：{item.players === null ? 0 : item.players.length}</span>
                                              <span>{item.status}</span>
                                          </Flex>
                                      }
                                      onClick={showDrawer1.bind(this, item)}
                                  />
                              </List.Item>)}
                    />
                    :
                    <div></div>
                }
            </div>
            <Drawer
                title="创建房间"
                placement={"bottom"}
                width={500}
                onClose={onClose}
                open={open}
                extra={
                    <Space>
                        <Button className="small-btn" onClick={onClose}>取消</Button>
                        <Button className="small-btn" type="primary" onClick={createRoom}>确定</Button>
                    </Space>
                }
            >
            </Drawer>
            <Drawer
                title="加入房间"
                placement={"bottom"}
                width={500}
                onClose={onClose1}
                open={open1}
                extra={
                    <Space>
                        <Button className="small-btn" onClick={onClose1}>取消</Button>
                        <Button className="small-btn" type="primary" onClick={joinRoom}>确定</Button>
                    </Space>
                }
            >
                <Space direction="vertical" size="middle">
                    <p>我的座位号</p>
                    <Space.Compact>
                        <Input type="text" className="input" value={seatPositionId} onChange={handsetSeatPosition}/>
                    </Space.Compact>
                    <p>我的名字</p>
                    <Space.Compact>
                        <Input className="input" value={playerName} onChange={handlePlayerNameChange}/>
                    </Space.Compact>
                </Space>
            </Drawer>
            <Context.Provider value={contextValue}>
                {contextHolder}
            </Context.Provider>
        </div>
    )
}

export default Home