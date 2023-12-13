function docReady(fn) {
    // see if DOM is already available
    if (document.readyState === "complete" || document.readyState === "interactive") {
        // call on next available tick
        setTimeout(fn, 1);
    } else {
        document.addEventListener("DOMContentLoaded", fn);
    }
}

class CheetahMessageEvent extends Event {
    constructor(type, data) {
        super(type);
        this.data = data;
    }
}

class CheetahWebsocket {
    constructor(url) {
        this.connected = false;
        this.url = url;
        this.source = new EventTarget();
        this._connect();
    }

    sendMessage(type, payload) {
        if (this.connected) {
            console.log(this.ws);
            this.ws.send(JSON.stringify({ type: type, payload: payload }));
        }
    }

    addListener(key, callback) {
        this.source.addEventListener(key, callback);
    }

    removeListener(key, callback) {
        this.source.removeEventListener(key, callback);
    }

    _connect() {
        console.log(`attempting to connect at ${this.url}`);
        const ws = new WebSocket(this.url);
        ws.onclose = this._onClose.bind(this);
        ws.onopen = this._onOpen.bind(this);
        ws.onmessage = this._onMessage.bind(this);
        ws.onerror = this._onError.bind(this);
        this.ws = ws;
    }

    _onMessage(e) {
        try {
            const data = JSON.parse(e.data);
            this.source.dispatchEvent(new CheetahMessageEvent(data.type, data.payload));
        } catch (e) {
            console.error(`[web-socket] onmessage error failed to parse ${e}`);
        }
    }

    _onError(e) {
        console.log(`[web-socket] error event: ${e}`);
    }
    
    _reconnectAttempt(count) {
        this._connect();
        setTimeout(() => {
            if(this.connected) {
                console.log("connected");
            }
        }, 1000);
    }

    _onClose(e) {
        console.log(`[web-socket] close event: ${e}`)
        this.connected = false;
        // retry to connect
        // this._reconnectAttempt(0);
    }

    _onOpen(e) {
        console.log(`[web-socket] open event: ${e}`);
        this.connected = true;
    }
}

docReady(() => {
    if (!window.LightweightCharts) {
        console.error("lightweight-charts is not isntalled.")
        return;
    }
    console.log("CREATING CHART")
    const { createChart } = window.LightweightCharts;
    const chartOptions = { layout: { textColor: 'black', background: { type: 'solid', color: 'white' } } };
    const chart = createChart(document.getElementById("container"), chartOptions);
    const candlestickSeries = chart.addCandlestickSeries({
        upColor: '#26a69a', downColor: '#ef5350', borderVisible: false,
        wickUpColor: '#26a69a', wickDownColor: '#ef5350',
    });
    var socket = new CheetahWebsocket("ws://localhost:8080/ws");
    socket.addListener("market:receive", (e) => {
        // console.log(e);
        let { start, open, high, low, close, volume } = e.data;
        open = Number.parseFloat(open);
        high = Number.parseFloat(high);
        low = Number.parseFloat(low);
        close = Number.parseFloat(close);
        candlestickSeries.update({ time: start, open: open, high: high, low: low, close: close });
    })
    const inpEle = document.getElementById("stock-search")
    const subBtn = document.getElementById("search-submit")
    subBtn.addEventListener("click", () => {
        console.log(inpEle.value);
        candlestickSeries.setData([]);
        socket.sendMessage("market:search", { symbol: inpEle.value });
    })
})