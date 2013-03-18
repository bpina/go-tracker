<h1>go-tracker</h1>
<span>A bittorrent tracker written in Go.</span>

<h2>Requirements</h2>
<ul>
  <li>A working Go installation.</li>
  <li>PostgreSQL</li>
</ul>

<h2>Installation</h2>
<pre>
git clone https://github.com/bpina/go-tracker.git
</pre>
<p>Edit the value for PREFIX in the makefile to your liking</p>
<pre>
make
make install clean
</pre>

<h2>Configuration</h2>
Add your database information to database.json in PREFIX/config.
<pre>
{
    "database": "database",
    "host": "localhost",
    "user": "user",
    "password": "password",
    "port": ""
}
</pre>

