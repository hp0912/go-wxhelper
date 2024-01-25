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

// 修改指令权限启用状态
function changeCommandEnableStatus(wxId) {
    axios({
        method: 'put',
        url: '/api/command/status',
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
        groupUsers.forEach((groupUser, i) => {
            console.log(groupUser)
            const { wxid, nickname, isMember, isAdmin, joinTime, lastActive, leaveTime, skipChatRank } = groupUser;

            let row = tbody.insertRow(i);
            // Insert data into cells
            row.insertCell(0).innerHTML = wxid;
            row.insertCell(1).innerHTML = nickname;
            row.insertCell(2).innerHTML = `<div class="badge badge-${isMember ? 'info' : 'error'} gap-2">${isMember ? '是' : '否'}</div>`;
            row.insertCell(3).innerHTML = `<div class="badge badge-${isAdmin ? 'info' : 'error'} gap-2">${isAdmin ? '是' : '否'}</div>`;
            row.insertCell(4).innerHTML = joinTime;
            row.insertCell(5).innerHTML = lastActive;
            row.insertCell(6).innerHTML = leaveTime;
            row.insertCell(7).innerHTML = `<input type="checkbox" class="toggle toggle-error" ${skipChatRank ? 'checked' : ''} onclick="changeUserGroupRankSkipStatus('${groupId}', '${wxid}')" />`;
        });
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
    }).finally(function () {
        // 隐藏加载框
        // loading.style.display = "none"
        groupNameTag.innerHTML = groupName
    })
}