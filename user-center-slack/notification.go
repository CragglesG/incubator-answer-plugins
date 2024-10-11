package slack_user_center

import (
	"strings"

	slackI18n "github.com/Anan1225/incubator-answer-plugins/user-center-slack/i18n"
	"github.com/apache/incubator-answer/plugin"
	"github.com/go-resty/resty/v2"
	"github.com/segmentfault/pacman/i18n"
	"github.com/segmentfault/pacman/log"
)

// GetNewQuestionSubscribers returns the subscribers of the new question notification
func (uc *UserCenter) GetNewQuestionSubscribers() (userIDs []string) {
	for userID, conf := range uc.UserConfigCache.userConfigMapping {
		if conf.AllNewQuestions {
			userIDs = append(userIDs, userID)
		}
	}
	return userIDs
}

// Notify sends a notification to the user
func (uc *UserCenter) Notify(msg plugin.NotificationMessage) {
	log.Debugf("try to send notification %+v", msg)

	if !uc.Config.Notification {
		return
	}

	// get user config
	userConfig, err := uc.getUserConfig(msg.ReceiverUserID)
	if err != nil {
		log.Errorf("get user config failed: %v", err)
		return
	}
	if userConfig == nil {
		log.Debugf("user %s has no config", msg.ReceiverUserID)
		return
	}

	// check if the notification is enabled
	switch msg.Type {
	case plugin.NotificationNewQuestion:
		if !userConfig.AllNewQuestions {
			log.Debugf("user %s not config the new question", msg.ReceiverUserID)
			return
		}
	case plugin.NotificationNewQuestionFollowedTag:
		if !userConfig.NewQuestionsForFollowingTags {
			log.Debugf("user %s not config the new question followed tag", msg.ReceiverUserID)
			return
		}
	case plugin.NotificationUpVotedTheAnswer:
		if !userConfig.UpvotedAnswers {
			log.Debugf("user %s not config the new upvoted answers", msg.ReceiverUserID)
		}
	case plugin.NotificationDownVotedTheAnswer:
		if !userConfig.DownvotedAnswers {
			log.Debugf("user %s not config the new downvoted answers", msg.ReceiverUserID)
		}

	case plugin.NotificationUpdateQuestion:
		if !userConfig.UpdatedQuestions {
			log.Debugf("user %s not config the update question", msg.ReceiverUserID)
			return
		}
	case plugin.NotificationUpdateAnswer:
		if !userConfig.UpdatedAnswers {
			log.Debugf("user %s not config the update answer", msg.ReceiverUserID)
			return
		}
	default:
		if !userConfig.InboxNotifications {
			log.Debugf("user %s not config the inbox notification", msg.ReceiverUserID)
			return
		}
	}

	log.Debugf("user %s config the notification", msg.ReceiverUserID)

	if len(userConfig.WebhookURL) == 0 {
		log.Errorf("user %s has no webhook url", msg.ReceiverUserID)
		return
	}

	notificationMsg := renderNotification(msg)
	// no need to send empty message
	if len(notificationMsg) == 0 {
		log.Debugf("this type of notification will be drop, the type is %s", msg.Type)
		return
	}

	// Create a Resty Client
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(NewWebhookReq(notificationMsg)).
		Post(userConfig.WebhookURL)

	if err != nil {
		log.Errorf("send message failed: %v %v", err, resp)
	} else {
		log.Infof("send message to %s success, resp: %s", msg.ReceiverUserID, resp.String())
	}
}

// renderNotification generates the notification message based on type
func renderNotification(msg plugin.NotificationMessage) string {
	lang := i18n.Language(msg.ReceiverLang)
	switch msg.Type {
	case plugin.NotificationUpdateQuestion:
		return plugin.TranslateWithData(lang, slackI18n.TplUpdatedQuestions, msg)
	case plugin.NotificationAnswerTheQuestion:
		return plugin.TranslateWithData(lang, slackI18n.TplAnswerTheQuestion, msg)
	case plugin.NotificationUpdateAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplUpdatedAnswers, msg)
	case plugin.NotificationAcceptAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplAcceptAnswer, msg)
	case plugin.NotificationCommentQuestion:
		return plugin.TranslateWithData(lang, slackI18n.TplCommentQuestion, msg)
	case plugin.NotificationCommentAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplCommentAnswer, msg)
	case plugin.NotificationReplyToYou:
		return plugin.TranslateWithData(lang, slackI18n.TplReplyToYou, msg)
	case plugin.NotificationMentionYou:
		return plugin.TranslateWithData(lang, slackI18n.TplMentionYou, msg)
	case plugin.NotificationInvitedYouToAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplInvitedYouToAnswer, msg)
	case plugin.NotificationNewQuestion, plugin.NotificationNewQuestionFollowedTag:
		msg.QuestionTags = strings.Join(strings.Split(msg.QuestionTags, ","), ", ")
		return plugin.TranslateWithData(lang, slackI18n.TplNewQuestion, msg)
	case plugin.NotificationUpVotedTheAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplUpvotedAnswer, msg)
	case plugin.NotificationDownVotedTheAnswer:
		return plugin.TranslateWithData(lang, slackI18n.TplDownvotedAnswer, msg)
	}
	return ""
}