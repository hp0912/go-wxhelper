console.log("打开首页")

//	改变AI开启状态
function changeAiEnableStatus(wxId) {
    // console.log("修改AI开启状态: ", wxId)

    axios({
        method: 'put',
        url: '/api/ai/status',
        data: {
            wxId: wxId
        }
    }).then(function (response) {
        console.log(`返回结果: ${JSON.stringify(response)}`);
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
        alert("修改失败")
    })
}

// 修改水群排行榜状态
function changeGroupRankEnableStatus(wxId) {
    // console.log("修改水群排行榜开启状态: ", wxId)
    axios({
        method: 'put',
        url: '/api/grouprank/status',
        data: {
            wxId: wxId
        }
    }).then(function (response) {
        console.log(`返回结果: ${JSON.stringify(response)}`);
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
        alert("修改失败")
    })
}

// 修改欢迎语开启状态
function changeWelcomeEnableStatus(wxId) {
    axios({
        method: 'put',
        url: '/api/welcome/status',
        data: {
            wxId: wxId
        }
    }).then(function (response) {
        console.log(`返回结果: ${JSON.stringify(response)}`);
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
        alert("修改失败")
    })
}

// 修改群成员是否参与排行榜状态
function changeUserGroupRankSkipStatus(groupId, userId) {
    console.log("修改水群排行榜开启状态: ", groupId, userId)
    axios({
        method: 'put',
        url: '/api/grouprank/skip',
        data: {
            wxId: groupId,
            userId: userId
        }
    }).then(function (response) {
        console.log(`返回结果: ${JSON.stringify(response)}`);
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
        alert("修改失败")
    })
}

// 获取群成员列表
function getGroupUsers(groupId, groupName) {
    // 获取表格的tbody部分，以便稍后向其中添加行
    var tbody = document.getElementById("groupUsers");
    tbody.innerHTML = ""

    // 打开模态框
    const modal = document.getElementById("groupUserModal");
    modal.showModal()

    // 设置群名称
    const groupNameTag = document.getElementById("groupUserModalName");
    groupNameTag.innerHTML = '<span class="loading loading-dots loading-lg"></span>'

    // 显示加载框
    // const loading = document.getElementById("groupUserDataLoading");
    // loading.style.display = "block"

    axios.get('/api/group/users', {
        params: {
            groupId: groupId
        }
    }).then(function (response) {
        // console.log(`返回结果: ${JSON.stringify(response)}`);
        // 渲染群成员列表
        const groupUsers = response.data
        // 循环渲染数据
        for (let i = 0; i < groupUsers.length; i++) {
            const groupUser = groupUsers[i]

            let row = tbody.insertRow(i); // 插入新行

            // 微信Id
            let wxId = row.insertCell(0);
            wxId.innerHTML = groupUser.wxid;

            // 昵称
            let nickname = row.insertCell(1);
            nickname.innerHTML = groupUser.nickname;

            // 是否群成员
            let isMember = row.insertCell(2);
            if (groupUser.isMember) {
                isMember.innerHTML = '<div class="badge badge-info gap-2">是</div>';
            } else {
                isMember.innerHTML = '<div class="badge badge-error gap-2">否</div>';
            }

            // 加群时间
            let joinTime = row.insertCell(3);
            joinTime.innerHTML = groupUser.joinTime;

            // 最后活跃时间
            let lastActiveTime = row.insertCell(4);
            lastActiveTime.innerHTML = groupUser.lastActiveTime;

            // 退群时间
            let leaveTime = row.insertCell(5);
            leaveTime.innerHTML = groupUser.leaveTime;

            // 是否跳过水群排行榜
            let skipChatRank = row.insertCell(6);
            skipChatRank.innerHTML = `<input type="checkbox" class="toggle toggle-error" ${groupUser.skipChatRank ? 'checked' : ''} onclick="changeUserGroupRankSkipStatus(\'${groupId}\', \'${groupUser.wxid}\')" />`;
        }
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
    }).finally(function () {
        // 隐藏加载框
        // loading.style.display = "none"
        groupNameTag.innerHTML = groupName
    })
}