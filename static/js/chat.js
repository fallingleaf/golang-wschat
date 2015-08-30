$(function() {

    var conn;
    var channellist = {};
    
    var chat = {type_chat: 0, type_method: 1, type_push: 2};
    chat.send_message = function(topic, msg) {
        var data = {
            type: this.type_chat,
            data: {
                topic: topic,
                message: msg
            }
        };
        conn.send(JSON.stringify(data));
        return false;
    };
    
    chat.call_method = function(params) {
        var data = {
            type: this.type_method,
            data: params
        };
        conn.send(JSON.stringify(data));
        return false;
    };
    
    var message = function(msg) {
        return {
            text: msg,
            status: 'unread'
        };
    };
    
    var channel = function(id, name) {
        var msgs = [];
        var lastmsg = 0;
        var unread = 0;
        var addMessage = function(msg) {
            this.msgs.push(msg);
        };
        
        var render = function(elem) {
            $(elem).empty();
            for(var i = 0; i < this.msgs.length; i++) {
                $(elem).append("<p class='msg'>"+this.msgs[i].text + "</p>");
            }
        };
        
        var addUnread = function() {
            this.unread += 1;
        };
        return {
            id: id,
            name: name,
            unread: unread,
            msgs: msgs,
            addUnread: addUnread,
            addMessage: addMessage,
            render: render
        };
        
    };
    
    $('.room-form').submit(function(evt) {
        evt.preventDefault();
        var room = $('#room-name').val();
        if(room != '') {
            var data = {
                api: "publish",
                topic: room
            };
            chat.call_method(data);
        }
        $('#room-name').val('');
        return false;
    });
    
    $('#form').submit(function(evt) {
        evt.preventDefault();
        var msg = $('#msg').val();
        if(!msg)
            return;
        var topic = $('li.active').attr('id');
        chat.send_message(topic, msg);
        $('#msg').val('');
        return false;
    });
    
    $('#topic-list').on('click', 'li', function(evt) {
        var id = $(this).attr('id');
        var name = $(this).find('.liname').text();
        $('.name').text(name);
        
        if(!channellist[id]) {
            var c = channel(id, name);
            channellist[id] = c;
            $(this).find('.label').remove();
            subscribeTopic(id);
        }
        
        $('#topic-list li').removeClass('active');
        $(this).addClass('active');
        
        var c = channellist[id];
        c.render('#log');
        c.unread = 0;
        $('#message-container').slideDown(); 
        return false;
    });
    
    function subscribeTopic(id) {
        var data = {
            api: "subscribe",
            topic: id
        };
        chat.call_method(data);
    }
    
    function newTopic(data) {
        if(channellist[data.id])
            return;
        var t = "<li id='" 
                + data.id 
                + "'><a href='#'><span class='liname'>"
                + data.name 
                + "</span><span class='badge'></span>"
                + "<span class='label label-warning'>subscribe</span>"
                + "</a></li>";
                
        $('#topic-list').append(t);
    }
    
    function handleMessage(evt) {
        var msg = JSON.parse(evt.data);
        if(!msg)
            return;
        var data = msg.data;
        //console.log(msg);
        switch(msg.type) {
            case chat.type_chat:
                var id = data.topic;
                var msg = data.message;
                var c = channellist[id];
                c.addMessage(message(msg));
                var selected = $('li.active').attr('id');
                if(id == selected) {
                    $('#log').append("<p class='msg'>" + msg + "</p");
                } else {
                    c.addUnread();
                    //$('li#'+id).addClass('new-msg');
                    $('li#'+id).find('span.badge').text(c.unread);
                }
                
                break;
            case chat.type_method:
                break;
            case chat.type_push:
                if(data.tag == 'publish') {
                    newTopic(data);
                }
                if(data.tag == 'sync') {
                    if(! data.topics)
                        return;
                    for(var i = 0; i < data.topics.length; i++) {
                        newTopic(data.topics[i]);
                    }
                }
                break;
            default:
                return false;
        }
    }
    
    function sync(evt) {
        var data = {
            api: 'syncTopic'
        };
        
        chat.call_method(data);
    }

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://go.tamhoangnguyen.me/ws");
        conn.onclose = function(evt) {
            //appendLog($("<div><b>Connection closed.</b></div>"))
        }
        conn.onmessage = handleMessage;
        conn.onopen = sync;
    } else {
        //appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
});