using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace CloudText
{
    public partial class ConfigForm : Form
    {
        string msg;
        public ConfigForm(string msg)
        {
            this.msg = msg;
            InitializeComponent();
        }

        private void Button1_Click(object sender, EventArgs e)
        {
            Properties.Settings.Default.Save();
            Application.Restart();
        }

        private void ConfigForm_Load(object sender, EventArgs e)
        {
            label5.Text = msg;
            if (msg == null || msg == "")
            {
                ControlBox = true;
            }
            else
            {
                ControlBox = false;
            }
        }
    }
}
