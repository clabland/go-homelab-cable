const scrollRate = 350;
const LCDWidth = 14;

function setCurrent(text) {
    const $current = document.getElementById("chan-current");
    $current.classList.remove("led-off");
    $current.classList.add("led-old");
    const interval = scrollLCD(text, $current);
    const clean = function() {
        clearInterval(interval);
        $current.classList.remove("led-old");
        $current.classList.add("led-off");
        $current.innerText = "88888888888888";
    }
    return clean;
}

function setNext(text) {
    const $next = document.getElementById("chan-next");
    $next.classList.remove("led-off");
    $next.classList.add("led-old");
    const interval = scrollLCD(text, $next);
    const clean = function() {
        clearInterval(interval);
        $next.classList.remove("led-old");
        $next.classList.add("led-off");
        $next.innerText = "88888888888888";
    }
    return clean;
}

function scrollLCD(text, $el) {
    let currIdx = 0;
    let direction = "forward";
    return setInterval(function () {
        // We can show whole thing
        if(text.length <= LCDWidth) {
            $el.innerText = text;
            return;
        }

        const display = text.substring(currIdx, currIdx+LCDWidth);

        // We've come back to the beginning
        if(currIdx == 0 && direction == "backward") {
            direction = "forward";
            $el.innerText = display;
            return;
        }

        // We need to go backwards
        if((currIdx + LCDWidth) == text.length && direction == "forward") {
            direction = "backward";
            $el.innerText = display;
            return;
        }

        if(direction == "forward") {
            currIdx++;
        } else {
            currIdx--;
        }

        $el.innerText = display;
    }, scrollRate);
}

const server = {
    endpoint: "http://localhost:3004/api",
    channels: [],
    selectedChannel: 0,
    channelLCDS: {
        clearCurrent: function(){},
        clearNext: function(){},
    },
    statusLCD: {
        $el: document.getElementById("status-lcd"),
        set: function(text, error) {
            this.$el.classList.remove("led-off");
            if(error) {
                this.$el.classList.remove("led-old");
                this.$el.classList.add("led-error");
                this.$el.innerText = text;
                return;
            }
            this.$el.classList.add("led-old");
            this.$el.innerText = text;
        },
        off: function() {
            this.$el.classList.remove("led-old");
            this.$el.classList.remove("led-error");
            this.$el.classList.add("led-off");
            this.$el.innerText = "88888888888888";
        }
    },
    powerOff: function() {
        server.channelLCDS.clearCurrent();
        server.channelLCDS.clearNext();
        for(var i=0; i<server.$chans.length; i++) {
            server.$chans[i].classList.remove("active");
            server.$chans[i].classList.add("disabled");
        }
        this.statusLCD.off();
    },
    error: function() {
        server.statusLCD.set("err - check console", true);
        throw new Error("halting for error");
    },
    enableChannels: function() {
        if(this.channels.length > server.$chans.length) {
            console.error(`can't understand more than ${server.$chans.length} channels.`)
            server.error();
            return;
        }
        for(var i=0; i<this.channels.length; i++) {
            server.$chans[i].classList.remove("disabled");
            if(i == 0) {
                server.channelLCDS.clearCurrent = setCurrent(this.channels[i].playing);
                server.channelLCDS.clearNext = setNext(this.channels[i].up_next);
            }
            if(this.channels[i].live){
                server.$chans[i].classList.add("active");
            }
        }
    },
    connect: async function() {
        server.$chans = document.querySelectorAll("div.channels > div > button");
        const channels = await server.getChannels();
        server.channels = channels;
        server.enableChannels();
        server.statusLCD.set("connected", false);
    },
    getChannels: async function() {
        const response = await fetch(`${server.endpoint}/networks`);
        const jsonData = await response.json();
        const channels = await fetch(`${server.endpoint}/networks/${jsonData[0].call_sign}/channels`);
        return Promise.resolve(channels.json());
    },
};

window.onload = async function() {
    const powerButton = document.querySelector("#power-btn > input[type=checkbox]");
    powerButton.addEventListener('change', async function() {
        try {
            if (this.checked) {
                await server.connect();
            } else {
                server.powerOff();
            }
        } catch(e) {
            console.log(e);
            server.error();
        }
    });
}