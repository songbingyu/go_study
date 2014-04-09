package main
 
import (
    "net/http"
    "fmt"
    "crypto/sha1"
    "io"
    "io/ioutil" 
    "sort"
    "encoding/xml"  
    "log" 
    "time"
    "bytes"
    "strconv"
    "x2j"
)

const (
    TEXT     = "text"
    IMAGE    = "image"
    LOCATION = "location"
    LINK     = "link"
    EVENT    = "event"
    VOICE    = "voice"
)


var  token = "weixinCourse";

func str2sha1(data string) string{
    t:=sha1.New()
    io.WriteString(t,data)
    return fmt.Sprintf("%x",t.Sum(nil))
}

func checkSignature( signature string,  timestamp string, nonce string )  bool {
	
	tmps:=[]string{token,timestamp,nonce}
	sort.Strings(tmps)
	tmpStr:=tmps[0]+tmps[1]+tmps[2]
    tmp:=str2sha1(tmpStr)
   if tmp==signature {
   		return true
   }

   return false
}

type IMessageSend interface{
    Text(Message)       Replay
    Image(Message)      Replay
    Location(Message)   Replay
    Link(Message)       Replay
    Event(Message)      Replay
    Voice(Message)      Replay
    Default(Message)    Replay
}


type MessageSendHelper struct {
   
}

func (w * MessageSendHelper) Text(msg *Message) (reply  Replay) {
    reply = Replay{} 
    reply.SetContent( string("OK"+msg.Content()) )
    return
}

func (w * MessageSendHelper) Image(msg *Message) (reply  Replay) {
  
    return
}
func (w * MessageSendHelper) Location(msg *Message) (reply  Replay) {
    reply = Replay{}
    log.Println("Mark : " + msg.Label())
    reply.SetContent("Mark : " + msg.Label())
    return
    return
}

func (w * MessageSendHelper) Link(msg *Message) (reply  Replay) {
    
    return
}

func (w * MessageSendHelper) Event(msg *Message) (reply  Replay) {
   
    return
}

func (w * MessageSendHelper) Voice(msg *Message) (reply  Replay) {
   
    return
}
func (w * MessageSendHelper) Default(msg *Message) (reply  Replay) {
   
    return
}


func MessageProc ( w http.ResponseWriter,msg *Message ) {


    var  reply Replay
    var  messagesendHelper MessageSendHelper
    var  ok  bool
    msgType:= msg.MsgType()
    switch msgType {
    case TEXT:
        reply = messagesendHelper.Text(msg)
    case IMAGE:
        reply = messagesendHelper.Image(msg)
    case LOCATION:
        reply = messagesendHelper.Location(msg)
    case LINK:
        reply = messagesendHelper.Link(msg)
    case EVENT:
        reply = messagesendHelper.Event(msg)
    case VOICE:
        reply = messagesendHelper.Voice(msg)
    default:
        reply = messagesendHelper.Default(msg)
    }
    if reply == nil {
        ok = true
        return // http 200
    }

    // auto-fix
    if reply.FromUserName() == "" {
        reply.SetFromUserName( msg.ToUserName())
    }
    if reply.ToUserName() == "" {
        reply.SetToUserName( msg.FromUserName() )
    }
    if reply.MsgType() == "" {
        reply.SetMsgType(TEXT)
    }

    reply.SetCreateTime(time.Now().Unix())

    if _, ok = reply["FuncFlag"]; !ok {
        reply.SetFuncFlag(0)
    }

    w.Write([]byte("<xml>"))
    _re := MapToXmlString(reply)
    w.Write([]byte(_re))
    w.Write([]byte("</xml>"))
    ok = true
} 



func HandleReq( w http.ResponseWriter, req *http.Request ) {
    
    fmt.Println("parse POST")  
  
    defer req.Body.Close()  
  
    body, err := ioutil.ReadAll( req.Body )  
    if err != nil {  
        log.Fatal(err)  
        return  
    }  
  
    fmt.Println(string(body))  

  
    var msg     BaseMsg

    if err = xml.Unmarshal( body, &msg ); err != nil {    
        log.Fatal(err)  
        return 
    }  

    root, err := x2j.DocToMap(string(body))
    if err != nil {
        fmt.Println("Bad XML Req", err)
        return
    }
    message := Message(root["xml"].(map[string]interface{}))
    fmt.Println(message)
    MessageProc(  w ,&message )

    return  
}


func Handler(w http.ResponseWriter, req *http.Request) {
    
     if req.Method == "GET" { 
	     // 微信加密签名
	        var signature = req.FormValue("signature");
	     // 时间戳
	        var timestamp = req.FormValue("timestamp");
	     // 随机数
	        var nonce = req.FormValue("nonce");
	     // 随机字符串
	      	var echostr = req.FormValue("echostr");
	     	
	     	
	     	if checkSignature( signature, timestamp, nonce) {

	     		fmt.Println( "signature=%s, timestamp=%s, nonce=%s,echostr=%s" , 
	       				 							signature, timestamp,nonce,echostr  );
	     		byteArray := []byte(echostr)
	     		w.Write( byteArray )
	     	} 
         
    } else if req.Method == "POST" {
    	 
       HandleReq( w, req )
    }  

}
 

func main() {

    http.HandleFunc("/monitor", Handler)
    http.ListenAndServe(":80", nil)
 
}

type  BaseMeszsage struct {
    attr string
}

type Message map[string]interface{}


type BaseMsg struct {

    XMLName xml.Name `xml:"xml"`
    field []string `xml:"any"`
}


//----------------------------------------------
func (w Message) String(key string) string {
    if str, ok := w[key]; ok {
        return str.(string)
    }
    return ""
}

func (w Message) Int64(key string) int64 {
    if val, ok := w[key]; ok {
        switch val.(type) {
        case string:
             i, _ := strconv.ParseInt(val.(string), 0, 64)
            return i
        case int:
            return int64(val.(int))
        case int64:
            return val.(int64)
        }
    }
    return 0
}


func (w Message) ToUserName() string {
    return w.String("ToUserName")
}
func (w Message) FromUserName() string {
    return w.String("FromUserName")
}
func (w Message) CreateTime() int64 {
    return w.Int64("CreateTime")
}
func (w Message) MsgType() string {
    return w.String("MsgType")
}
func (w Message) MsgId() string {
    return w.String("MsgId")
}

//------------------------------------

func (w Message) Content() string {
    return w.String("Content")
}

//-----------------------------------

func (w Message) PicUrl() string {
    return w.String("PicUrl")
}

//-----------------------------------

func (w Message) Location_X() string {
    return w.String("Location_X")
}
func (w Message) Location_Y() string {
    return w.String("Location_Y")
}
func (w Message) Scale() int64 {
    return w.Int64("Scale")
}
func (w Message) Label() string {
    return w.String("Label")
}

//--------------------------------

func (w Message) Event() string {
    return w.String("Event")
}
func (w Message) EventKey() string {
    return w.String("EventKey")
}

//--------------------------------
func (w Message) Title() string {
    return w.String("Title")
}
func (w Message) Description() string {
    return w.String("Description")
}
func (w Message) Url() string {
    return w.String("Url")
}

//-------------------------------
func (w Message) MediaId() string {
    return w.String("MediaId")
}

func (w Message) Format() string {
    return w.String("Format")
}




type Replay map[string]interface{}

func (r Replay) String(key string) string {
    if str, ok := r[key]; ok {
        return str.(string)
    }
    return ""
}

func (r Replay) Int64(key string) int64 {
    if val, ok := r[key]; ok {
        switch val.(type) {
        case string:
            i, _ := strconv.ParseInt(val.(string), 0, 64)
            return i
        case int:
            return int64(val.(int))
        case int64:
            return val.(int64)
        }
    }
    return 0
}

func (r Replay) ToUserName() string {
    return r.String("ToUserName")
}
func (r Replay) FromUserName() string {
    return r.String("FromUserName")
}
func (r Replay) CreateTime() int64 {
    return r.Int64("CreateTime")
}
func (r Replay) MsgType() string {
    return r.String("MsgType")
}
func (r Replay) FuncFlag() int64 {
    return r.Int64("FuncFlag")
}

func (r Replay) SetToUserName(val string) Replay {
    r["ToUserName"] = val
    return r
}
func (r Replay) SetFromUserName(val string) Replay {
    r["FromUserName"] = val
    return r
}
func (r Replay) SetCreateTime(val int64) Replay {
    r["CreateTime"] = val
    return r
}
func (r Replay) SetMsgType(val string) Replay {
    r["MsgType"] = val
    return r
}
func (r Replay) SetFuncFlag(val int64) Replay {
    r["FuncFlag"] = val
    return r
}

//----------------------------------------
func (r Replay) Content() string {
    return r.String("Content")
}
func (r Replay) SetContent(val string) Replay {
    r["Content"] = val
    return r
}

//----------------------------------------

type MusicOut struct {
    ToUserName   string
    FromUserName string
    CreateTime   int64
    MsgType      string

    Title       string `xml:"Music>Title"`
    Description string `xml:"Music>Description"`
    MusicUrl    string `xml:"Music>MusicUrl"`
    HQMusicUrl  string `xml:"Music>HQMusicUrl"`
    FuncFlag    int
}

func MapToXmlString(m map[string]interface{}) string {
    buf := &bytes.Buffer{}
    for k, v := range m {

        if v != nil {
            switch v.(type) {
            case int:
                io.WriteString(buf, fmt.Sprintf("<%s>", k))
                io.WriteString(buf, fmt.Sprintf("%d", v))
                io.WriteString(buf, fmt.Sprintf("</%s>\n", k))
            case int64:
                io.WriteString(buf, fmt.Sprintf("<%s>", k))
                io.WriteString(buf, fmt.Sprintf("%d", v))
                io.WriteString(buf, fmt.Sprintf("</%s>\n", k))
            case string:
                io.WriteString(buf, fmt.Sprintf("<%s>", k))
                io.WriteString(buf, "<![CDATA["+v.(string)+"]]>")
                io.WriteString(buf, fmt.Sprintf("</%s>\n", k))
            case map[string]interface{}:
                io.WriteString(buf, fmt.Sprintf("<%s>", k))
                io.WriteString(buf, MapToXmlString(v.(map[string]interface{})))
                io.WriteString(buf, fmt.Sprintf("</%s>\n", k))
            case []interface{}:
                for _, t := range v.([]interface{}) {
                    switch t.(type) {
                    case map[string]interface{}:
                        io.WriteString(buf, fmt.Sprintf("<%s>", k))
                        io.WriteString(buf, MapToXmlString(t.(map[string]interface{})))
                        io.WriteString(buf, fmt.Sprintf("</%s>\n", k))
                    }
                }
            }
        }

    }
    return buf.String()
}