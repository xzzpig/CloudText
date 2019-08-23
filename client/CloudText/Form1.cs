using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Windows.Forms;
using Websocket.Client;

namespace CloudText
{

    public partial class Form1 : Form
    {
        [DllImport("user32.dll", CharSet = CharSet.Auto)]
        public static extern int AddClipboardFormatListener(IntPtr hWnd);
        [DllImport("user32.dll", CharSet = CharSet.Auto)]
        public static extern int RemoveClipboardFormatListener(IntPtr hWnd);

        private WebsocketClient wsclient;

        public Form1()
        {
            InitializeComponent();
        }

        private void Form1_Load(object sender, EventArgs e)
        {
            notifyIcon1.Visible = true;
            ShowInTaskbar = false;
            Visible = false;
            AddClipboardFormatListener(Handle);
            connectWs();
        }
        private async void connectWs()
        {
            Console.WriteLine(Properties.Settings.Default.wsuri);
            var ws = new WebsocketClient(new Uri(Properties.Settings.Default.wsuri));
            _ = Task.Run(async () =>
            {
                await Task.Delay(1000);
                if (!ws.IsRunning)
                {
                    new ConfigForm("无法连接到websocket").ShowDialog();
                }
            });

            ws.MessageReceived.Subscribe(action => {
                var message = action.Text;
                Console.WriteLine(message);
                var response = JsonConvert.DeserializeObject<CloudTextPackage>(message);
                switch (response.action)
                {
                    case "Auth":
                        if (response.data != "success")
                        {
                            new ConfigForm("账号登录失败").ShowDialog();
                        }
                        else
                        {
                            this.wsclient = ws;
                        }
                        break;
                    case "Set":
                        var t = new Thread(() => {
                            try
                            {
                                if (Clipboard.GetText() != response.data&&Properties.Settings.Default.autoget==CheckState.Checked)
                                    Clipboard.SetText(response.data);
                            }
                            catch (Exception e) { }
                        });
                        t.SetApartmentState(ApartmentState.STA);
                        t.Start();
                        break;
                }

            });
            var req = new CloudTextPackage(0, Properties.Settings.Default.username, Properties.Settings.Default.password, "Auth", null);
            await ws.Start();
            await ws.Send(JsonConvert.SerializeObject(req));
        }
        private void 设置ToolStripMenuItem_Click(object sender, EventArgs e)
        {
            new ConfigForm("").ShowDialog();
        }
        protected override void WndProc(ref Message m)
        {
            base.WndProc(ref m);
            Console.WriteLine(m);
            if (m.Msg == 0x31d && wsclient != null && Properties.Settings.Default.autoset == CheckState.Checked)
            {
                var req = new CloudTextPackage(0, null, null, "Set", Clipboard.GetText());
                wsclient.Send(JsonConvert.SerializeObject(req));
            }
        }

        private void 退出ToolStripMenuItem_Click(object sender, EventArgs e)
        {
            Close();
        }

        private void autoget_ToolStripMenuItem_CheckStateChanged(object sender, EventArgs e)
        {
            Properties.Settings.Default.autoget = ((ToolStripMenuItem)sender).CheckState;
            Properties.Settings.Default.Save();
            Application.Restart();
        }

        private void autoset_ToolStripMenuItem_CheckStateChanged(object sender, EventArgs e)
        {
            Properties.Settings.Default.autoset = ((ToolStripMenuItem)sender).CheckState;
            Properties.Settings.Default.Save();
            Application.Restart();
        }

        private void Form1_FormClosing(object sender, FormClosingEventArgs e)
        {
            RemoveClipboardFormatListener(this.Handle);
        }

        private void NotifyIcon1_MouseDoubleClick(object sender, MouseEventArgs e)
        {
            new ConfigForm("").ShowDialog();
        }
    }
}
