import QtQuick 2.4
import Ubuntu.Components 1.3
import "../../components"
import Ubuntu.Content 1.1

TelegramPage {

    id: addlinkdevice
    width: units.gu(40)
    title: i18n.tr("Link a Device")
    height: units.gu(68)
    visible: true
    // header: PageHeader {
    //   title: i18n.tr("Link a Device")
    //     id: pageHeader
    //     clip:true
    //     width: parent.width
    //     height: units.gu(6)
    // }

    Rectangle {
        width: parent.width -100
        height: parent.height -100
        color: "red"
        border.color: "black"
        border.width: 5
        radius: 10
        anchors.verticalCenter: parent.verticalCenter
        anchors.horizontalCenter: parent.horizontalCenter
        Camera {

        }
    }


}
