/**
 * Form.js 表单提交类
 * 依赖 jquery.form.js
 */

var Form = {

    /**
     * 失败 div
     */
    failedBox: '#failedBox',

    /**
     * 是否在弹框中
     */
    inPopup: false,

    /**
     * ajax submit
     * @param element
     * @param inPopup
     * @returns {boolean}
     */
    ajaxSubmit: function(element, inPopup) {

        if (inPopup) {
            Form.inPopup = true;
        }

        /**
         * 成功弹出信息
         * @param message
         * @param data
         */
        function successPopup(message, data) {
            Layers.success(message)
        }

        /**
         * 成功信息条
         * @param message
         * @param data
         */
        function successBox(message, data) {
            var text = [message];

            $(Form.failedBox).html('');
            $(Form.failedBox).removeClass();
            $(Form.failedBox).addClass('alert alert-success');
            $(Form.failedBox).append('<a class="close" href="#" onclick="$(this).parent().hide();">×</a>');
            $(Form.failedBox).append('<strong >操作成功！ </strong>');
            var ul = $('<ul></ul>');
            for (var i = 0; i < text.length; i++) {
                ul.append('<li>' + text[i] + '</li>');
            }
            $(Form.failedBox).append(ul);
            $(Form.failedBox).show();
        }

        /**
         * 失败信息条
         * @param message
         * @param data
         */
        function failed(message, data) {
            var text = [message];

            $(Form.failedBox).html('');
            $(Form.failedBox).removeClass('hide');
            $(Form.failedBox).addClass('alert alert-danger');
            $(Form.failedBox).append('<a class="close" href="#" onclick="$(this).parent().hide();">×</a>');
            $(Form.failedBox).append('<strong >操作失败！ </strong>');
            var ul = $('<ul></ul>');
            for (var i = 0; i < text.length; i++) {
                ul.append('<li>' + text[i] + '</li>');
            }
            $(Form.failedBox).append(ul);
            $(Form.failedBox).show();
        }

        /**
         * response
         * @param result
         */
        function response(result) {
            //console.log(result)
            if (result.code == 0) {
                failed(result.message, result.data);
            }
            if (result.code == 1) {
                successBox(result.message, result.data);
            }
            if (result.redirect.url) {
                var sleepTime = result.redirect.sleep || 3000;
                setTimeout(function() {
                    if (Form.inPopup) {
                        parent.location.href = result.redirect.url;
                    } else {
                        location.href = result.redirect.url;
                    }
                }, sleepTime);
            }
        }

        var options = {
            dataType: 'json',
            success: response
        };

        $(element).ajaxSubmit(options);

        return false;
    }
};