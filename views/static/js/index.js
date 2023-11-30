console.log("打开首页")

//	改变AI开启状态
function changeAiEnableStatus(wxId) {
    console.log("修改AI开启状态: ", wxId)

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
    })
}

// 修改水群排行榜状态
function changeGroupRankEnableStatus(wxId) {
    console.log("修改水群排行榜开启状态: ", wxId)
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
    })
}

// 修改群成员是否参与排行榜状态
function changeUserGroupRankSkipStatus(groupId, userId) {
    console.log("修改水群排行榜开启状态: ", groupId, userId)
    axios({
        method: 'put',
        url: '/api/grouprank/skip',
        data: {
            wxId: wxId,
            userId: userId
        }
    }).then(function (response) {
        console.log(`返回结果: ${JSON.stringify(response)}`);
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
    })
}

// 获取群成员列表
function getGroupUsers(groupId) {
    // 打开模态框
    const modal = document.getElementById("groupUserModal");
    modal.showModal()

    axios.get('/api/group/users', {
        params: {
            groupId: groupId
        }
    }).then(function (response) {
        console.log(`返回结果: ${JSON.stringify(response)}`);
        // 渲染群成员列表
        const groupUsers = response.data

        // const groupUserList = document.getElementById("groupUsers")

        // 获取表格的tbody部分，以便稍后向其中添加行
        var tbody = document.getElementById("groupUsers");
        for (let i = 0; i < groupUsers.length; i++) {
            const groupUser = groupUsers[i]

            var row = tbody.insertRow(i); // 插入新行

            // 微信Id
            var wxId = row.insertCell(0);
            wxId.innerHTML = data[i].wxId;

            // 昵称
            var nickname = row.insertCell(1);
            nickname.innerHTML = data[i].wxId;

            // 是否群成员
            var isMember = row.insertCell(2);
            isMember.innerHTML = data[i].wxId;

            // 最后活跃时间
            var wxId = row.insertCell(3);
            wxId.innerHTML = data[i].wxId;

            // 退群时间
            var wxId = row.insertCell(4);
            wxId.innerHTML = data[i].wxId;

            // 是否跳过水群排行榜
            var wxId = row.insertCell(5);
            wxId.innerHTML = data[i].wxId;
        }
    }).catch(function (error) {
        console.log(`错误信息: ${error}`);
    })
}