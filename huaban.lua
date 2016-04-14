local cs = "hello lua"
local dir = "/root/Doc/bin/huaban"
-- local dir = "warms"

function makecommand( )
	local one = "cd " .. dir
	local two = "git pull"
	return one .. ";" .. two
end

function gitpull( )
	local cmd = makecommand()
	print(cmd)
	os.execute(cmd)
end

gitpull()