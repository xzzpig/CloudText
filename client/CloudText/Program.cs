using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace CloudText
{
    public class CloudTextPackage
    {
        public int id;
        public string username;
        public string password;
        public string action;
        public string data;

        public CloudTextPackage(int id, string username, string password, string action, string data)
        {
            this.id = id;
            this.username = username;
            this.password = password;
            this.action = action;
            this.data = data;
        }
    }
    static class Program
    {
        /// <summary>
        /// 应用程序的主入口点。
        /// </summary>
        [STAThread]
        static void Main()
        {
            Application.EnableVisualStyles();
            Application.SetCompatibleTextRenderingDefault(false);
            Application.Run(new Form1());
        }
    }
}
